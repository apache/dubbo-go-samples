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
package ollama

import (
	"github.com/apache/dubbo-go-samples/book-flight-ai-agent/go-server/model"
)

func WithStreamingFunc(fn model.CallFunc) model.Option {
	return &fn
}

func WithNumberKeep(numKeep int64) model.Option {
	return map[string]any{"num_keep": numKeep}
}

func WithSeed(seed int64) model.Option {
	return map[string]any{"seed": seed}
}

func WithNumProdict(num_predict int64) model.Option {
	return map[string]any{"num_predict": num_predict}
}

func WithTopK(topK int64) model.Option {
	return map[string]any{"top_k": topK}
}

func WithTopP(top_p float64) model.Option {
	return map[string]any{"top_p": top_p}
}

func WithMinP(min_p float64) model.Option {
	return map[string]any{"min_p": min_p}
}

func WithTypicalP(typical_p float64) model.Option {
	return map[string]any{"typical_p": typical_p}
}

func WithRepeatLastN(repeat_last_n int64) model.Option {
	return map[string]any{"repeat_last_n": repeat_last_n}
}

func WithTemperature(temperature float64) model.Option {
	return map[string]any{"temperature": temperature}
}

func WithRepeatPenalty(repeat_penalty float64) model.Option {
	return map[string]any{"repeat_penalty": repeat_penalty}
}

func WithPresencePenalty(presence_penalty float64) model.Option {
	return map[string]any{"presence_penalty": presence_penalty}
}

func WithFrequencyPenalty(frequency_penalty float64) model.Option {
	return map[string]any{"frequency_penalty": frequency_penalty}
}

func WithMirostat(mirostat int64) model.Option {
	return map[string]any{"mirostat": mirostat}
}

func WithMirostatTau(mirostat_tau float64) model.Option {
	return map[string]any{"mirostat_tau": mirostat_tau}
}

func WithMirostatEta(mirostat_eta float64) model.Option {
	return map[string]any{"mirostat_eta": mirostat_eta}
}

func WithPenalizeNewline(penalize_newline bool) model.Option {
	return map[string]any{"penalize_newline": penalize_newline}
}

func WithStop(stop []string) model.Option {
	return map[string]any{"stop": stop}
}

func WithNama(nama bool) model.Option {
	return map[string]any{"nama": nama}
}

// Default: 1024
func WithNumberContext(numCtx int64) model.Option {
	return map[string]any{"num_ctx": numCtx}
}

func WithNumberBatch(num_batch int64) model.Option {
	return map[string]any{"num_batch": num_batch}
}

func WithNumberGPU(num_gpu int64) model.Option {
	return map[string]any{"num_gpu": num_gpu}
}

func WithMainGPU(main_gpu int64) model.Option {
	return map[string]any{"main_gpu": main_gpu}
}

func WithLowVram(low_vram bool) model.Option {
	return map[string]any{"low_vram": low_vram}
}

func WithVocabOnly(vocab_only bool) model.Option {
	return map[string]any{"vocab_only": vocab_only}
}

func WithVuseMmap(vuse_mmap bool) model.Option {
	return map[string]any{"vuse_mmap": vuse_mmap}
}

func WithUseMLock(use_mlock bool) model.Option {
	return map[string]any{"use_mlock": use_mlock}
}

func WithNumberThread(num_thread int64) model.Option {
	return map[string]any{"num_thread": num_thread}
}
