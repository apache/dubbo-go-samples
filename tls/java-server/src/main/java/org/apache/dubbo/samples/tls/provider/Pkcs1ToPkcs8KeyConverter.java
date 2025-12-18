/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.dubbo.samples.tls.provider;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.attribute.PosixFilePermission;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Base64;
import java.util.HashSet;
import java.util.Set;
import java.util.ArrayList;
import java.util.List;

/**
 * Converts PKCS#1 format private keys (-----BEGIN RSA PRIVATE KEY-----)
 * to PKCS#8 format (-----BEGIN PRIVATE KEY-----) that Java/Netty can parse.
 * 
 * Uses openssl pkcs8 command internally via system process.
 * Implements secure temporary file handling with proper cleanup and permissions.
 */
public class Pkcs1ToPkcs8KeyConverter {
    
    // Track temporary files for cleanup
    private static final List<File> tempFiles = new ArrayList<>();
    
    // Install shutdown hook once
    static {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            cleanupTempFiles();
        }));
    }

    /**
     * Convert PKCS#1 private key to PKCS#8 format
     * @param pkcs1KeyPath path to PKCS#1 PEM file
     * @return PKCS#8 PEM content as string
     * @throws Exception if conversion fails
     */
    public static String convertPkcs1ToPkcs8(String pkcs1KeyPath) throws Exception {
        // Validate input path
        if (pkcs1KeyPath == null || pkcs1KeyPath.trim().isEmpty()) {
            throw new IllegalArgumentException("Key path cannot be null or empty");
        }
        
        // Verify input file exists and is readable
        File inputFile = new File(pkcs1KeyPath);
        if (!inputFile.exists()) {
            throw new FileNotFoundException("Private key file not found: " + pkcs1KeyPath);
        }
        if (!inputFile.canRead()) {
            throw new IOException("Cannot read private key file: " + pkcs1KeyPath);
        }
        
        // Call openssl to convert PKCS#1 to PKCS#8
        ProcessBuilder pb = new ProcessBuilder(
                "openssl", "pkcs8",
                "-topk8",
                "-inform", "PEM",
                "-outform", "PEM",
                "-in", pkcs1KeyPath,
                "-nocrypt"
        );
        
        Process process = null;
        try {
            process = pb.start();
        } catch (IOException e) {
            throw new RuntimeException("Failed to execute openssl command. Ensure OpenSSL is installed and in PATH.", e);
        }
        
        final Process finalProcess = process;
        final StringBuilder output = new StringBuilder();
        final StringBuilder errorOutput = new StringBuilder();
        
        try {
            // Read both stdout and stderr concurrently to prevent deadlock
            Thread outputThread = new Thread(() -> {
                try (BufferedReader reader = new BufferedReader(new InputStreamReader(finalProcess.getInputStream()))) {
                    String line;
                    while ((line = reader.readLine()) != null) {
                        synchronized (output) {
                            output.append(line).append("\n");
                        }
                    }
                } catch (IOException e) {
                    System.err.println("[TLS WARNING] Error reading process output: " + e.getMessage());
                }
            }, "openssl-output-reader");
            
            Thread errorThread = new Thread(() -> {
                try (BufferedReader reader = new BufferedReader(new InputStreamReader(finalProcess.getErrorStream()))) {
                    String line;
                    while ((line = reader.readLine()) != null) {
                        synchronized (errorOutput) {
                            errorOutput.append(line).append("\n");
                        }
                    }
                } catch (IOException e) {
                    System.err.println("[TLS WARNING] Error reading process error stream: " + e.getMessage());
                }
            }, "openssl-error-reader");
            
            outputThread.start();
            errorThread.start();
            
            // Wait for process to complete
            int exitCode = process.waitFor();
            
            // Wait for stream readers to finish
            outputThread.join(5000); // 5 second timeout
            errorThread.join(5000);
            
            if (exitCode != 0) {
                throw new RuntimeException("OpenSSL PKCS#1 to PKCS#8 conversion failed (exit code " + exitCode + "): " + errorOutput.toString());
            }
            
            if (output.length() == 0) {
                throw new RuntimeException("OpenSSL produced no output. Error: " + errorOutput.toString());
            }
            
            return output.toString();
            
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("OpenSSL process was interrupted", e);
        } finally {
            // Ensure process is destroyed to free resources
            if (process != null) {
                process.destroy();
                // Give it a moment to terminate gracefully
                try {
                    if (!process.waitFor(1000, java.util.concurrent.TimeUnit.MILLISECONDS)) {
                        process.destroyForcibly();
                    }
                } catch (InterruptedException e) {
                    process.destroyForcibly();
                    Thread.currentThread().interrupt();
                }
            }
        }
    }
    
    /**
     * Check if a PEM file is in PKCS#1 format
     * @param pemPath path to PEM file
     * @return true if it contains "BEGIN RSA PRIVATE KEY"
     * @throws IOException if file read fails
     */
    public static boolean isPkcs1Format(String pemPath) throws IOException {
        try (BufferedReader reader = new BufferedReader(new FileReader(pemPath))) {
            String firstLine = reader.readLine();
            return firstLine != null && firstLine.contains("BEGIN RSA PRIVATE KEY");
        }
    }
    
    /**
     * Write PKCS#8 PEM content to a temporary file with secure permissions
     * @param pkcs8Content PEM content
     * @return path to temporary file
     * @throws IOException if file creation fails
     */
    public static String writeTempPkcs8File(String pkcs8Content) throws IOException {
        File temp = File.createTempFile("pkcs8_key_", ".pem");
        
        // Set restrictive permissions (owner read/write only) on Unix-like systems
        try {
            Set<PosixFilePermission> perms = new HashSet<>();
            perms.add(PosixFilePermission.OWNER_READ);
            perms.add(PosixFilePermission.OWNER_WRITE);
            Files.setPosixFilePermissions(temp.toPath(), perms);
            System.out.println("[TLS] Set secure permissions (600) on temporary key file");
        } catch (UnsupportedOperationException e) {
            // Windows or non-POSIX system - file permissions not supported
            System.out.println("[TLS WARNING] Cannot set POSIX permissions on this system. Temporary key file may be accessible to other users.");
        }
        
        // Register for cleanup on exit
        temp.deleteOnExit();
        synchronized (tempFiles) {
            tempFiles.add(temp);
        }
        
        try (FileWriter writer = new FileWriter(temp)) {
            writer.write(pkcs8Content);
        }
        
        return temp.getAbsolutePath();
    }
    
    /**
     * Clean up all temporary key files securely
     */
    public static void cleanupTempFiles() {
        synchronized (tempFiles) {
            for (File file : tempFiles) {
                if (file.exists()) {
                    try {
                        // Overwrite file content before deletion for security
                        secureDelete(file);
                        System.out.println("[TLS] Securely deleted temporary key file: " + file.getAbsolutePath());
                    } catch (Exception e) {
                        System.err.println("[TLS WARNING] Failed to securely delete temp file: " + file.getAbsolutePath() + ", " + e.getMessage());
                        // Attempt regular deletion as fallback
                        file.delete();
                    }
                }
            }
            tempFiles.clear();
        }
    }
    
    /**
     * Securely delete a file by overwriting its content before deletion
     * @param file file to delete
     * @throws IOException if deletion fails
     */
    private static void secureDelete(File file) throws IOException {
        if (file.exists()) {
            long length = file.length();
            try (RandomAccessFile raf = new RandomAccessFile(file, "rws")) {
                // Overwrite with zeros
                raf.seek(0);
                for (long i = 0; i < length; i++) {
                    raf.write(0);
                }
            }
            // Delete the file
            if (!file.delete()) {
                throw new IOException("Failed to delete file: " + file.getAbsolutePath());
            }
        }
    }
    
    /**
     * Load and convert PKCS#1 key if necessary
     * @param keyPath path to private key file
     * @return path to PKCS#8 format key (may be temporary file)
     * @throws Exception if conversion fails
     */
    public static String loadAndConvertKey(String keyPath) throws Exception {
        if (!isPkcs1Format(keyPath)) {
            // Already PKCS#8 or other format, use as is
            return keyPath;
        }
        
        System.out.println("[TLS] Detected PKCS#1 format key, converting to PKCS#8...");
        String pkcs8Content = convertPkcs1ToPkcs8(keyPath);
        String tempPath = writeTempPkcs8File(pkcs8Content);
        System.out.println("[TLS] Converted key written to temporary file: " + tempPath);
        System.out.println("[TLS] Temporary key file will be securely deleted on JVM shutdown");
        
        return tempPath;
    }
}
