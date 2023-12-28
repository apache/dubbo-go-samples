package org.apache.dubbo.server;

import org.apache.dubbo.config.spring.context.annotation.EnableDubbo;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * @author zhaoyunxing
 */
@EnableDubbo
@SpringBootApplication
public class JavaServerApplication {

	public static void main(String[] args) {
		SpringApplication.run(JavaServerApplication.class, args);
	}

}
