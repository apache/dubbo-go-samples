/**
 * Game main logic
 * Using Vue 3 Composition API
 */

const { createApp, ref, reactive, onMounted, onUnmounted } = Vue;
const { ElMessage } = ElementPlus;

createApp({
    setup() {
        // Reactive data
        const dialogVisible = ref(true);
        const info = reactive({
            name: '',
            score: 0,
            rank: 0,
        });
        const isMeet = ref(false);
        const loadingProgress = ref(0);

        // Animation interval IDs storage for cleanup
        const intervals = ref({
            progress: null,
            gate: null,
            ball: null,
        });

        /**
         * Get rank
         */
        const getRank = async () => {
            try {
                const response = await rank({ name: info.name });
                if (response.data && response.data.rank !== undefined) {
                    info.rank = response.data.rank;
                }
            } catch (error) {
                ElMessage({
                    message: `Failed to get rank: ${error.message || error}`,
                    type: 'error',
                });
            }
        };

        /**
         * Submit user info
         */
        const submitInfo = async () => {
            if (!info.name.trim()) {
                ElMessage({
                    message: 'Please enter your name',
                    type: 'warning',
                });
                return;
            }

            dialogVisible.value = false;
            try {
                const response = await login({ name: info.name });
                if (response.data && response.data.to) {
                    info.name = response.data.to;
                }
                if (response.data && response.data.score !== undefined) {
                    info.score = response.data.score;
                }
                ElMessage({
                    message: `Welcome ${response.msg || info.name}`,
                    type: 'success',
                });
                await getRank();
            } catch (error) {
                ElMessage({
                    message: `Login failed: ${error.message || error}`,
                    type: 'error',
                });
                dialogVisible.value = true;
            }
        };

        /**
         * Move ball (kick animation)
         */
        const moveBall = () => {
            const ballElement = document.querySelector('.ball');
            const gateElement = document.getElementById('gate');
            
            if (!ballElement || !gateElement) {
                return;
            }

            // Clear previous animation
            if (intervals.value.ball) {
                clearInterval(intervals.value.ball);
            }

            let tempHeight = 0;
            const clientHeight = document.documentElement.clientHeight;

            intervals.value.ball = setInterval(() => {
                const ballRect = ballElement.getBoundingClientRect();
                const gateRect = gateElement.getBoundingClientRect();
                
                const space = ballRect.left - gateRect.left;
                const height = ballRect.top - gateRect.top;

                // Detect goal (fixed collision detection logic)
                if (!isMeet.value && space < 60 && space > -20 && height < 20 && height > -20) {
                    isMeet.value = true;
                    clearInterval(intervals.value.ball);
                    intervals.value.ball = null;

                    // Submit score
                    score({ name: info.name, score: 1 })
                        .then((response) => {
                            if (response.data && response.data.score !== undefined) {
                                info.score = response.data.score;
                            }
                            ElMessage({
                                message: `Goal! Total score: ${response.data.score || info.score}`,
                                type: 'success',
                            });
                            getRank();
                        })
                        .catch((error) => {
                            ElMessage({
                                message: `Failed to record score: ${error.message || error}`,
                                type: 'error',
                            });
                        });
                }

                // Ball moves to bottom
                if (tempHeight >= clientHeight) {
                    ballElement.style.bottom = '0px';
                    clearInterval(intervals.value.ball);
                    intervals.value.ball = null;
                } else {
                    tempHeight += 20;
                    ballElement.style.bottom = `${tempHeight}px`;
                }
            }, 20);
        };

        /**
         * Loading progress bar animation
         */
        const moveProgress = () => {
            const progressBar = document.getElementById('myBar');
            if (!progressBar) {
                return;
            }

            let width = 0;
            intervals.value.progress = setInterval(() => {
                if (width >= 100) {
                    clearInterval(intervals.value.progress);
                    intervals.value.progress = null;
                    progressBar.remove();
                } else {
                    width++;
                    loadingProgress.value = width;
                    progressBar.style.width = `${width}%`;
                    progressBar.textContent = `${width}%`;
                }
            }, 20);
        };

        /**
         * Move gate (left-right animation)
         */
        const moveGate = () => {
            const gateElement = document.getElementById('gate');
            if (!gateElement) {
                return;
            }

            // Get initial marginLeft (set to -4rem in CSS, which is -64px)
            const initialMargin = -64; // 8rem / 2 = 4rem = 64px
            let tempValue = 0;
            let increase = 1;

            intervals.value.gate = setInterval(() => {
                if (tempValue > 100 || tempValue < -100) {
                    increase = -increase;
                }
                tempValue += increase;
                // Move based on initial position
                gateElement.style.marginLeft = `${initialMargin + tempValue}px`;
            }, 20);
        };

        /**
         * Click person to kick ball
         */
        const clickPeople = () => {
            isMeet.value = false;
            moveBall();
        };

        // Start animations when component is mounted
        onMounted(() => {
            moveProgress();
            moveGate();
        });

        // Clean up all animations when component is unmounted
        onUnmounted(() => {
            if (intervals.value.progress) {
                clearInterval(intervals.value.progress);
            }
            if (intervals.value.gate) {
                clearInterval(intervals.value.gate);
            }
            if (intervals.value.ball) {
                clearInterval(intervals.value.ball);
            }
        });

        return {
            dialogVisible,
            info,
            isMeet,
            loadingProgress,
            submitInfo,
            clickPeople,
        };
    },
}).use(ElementPlus).mount('#app');
