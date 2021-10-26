new Vue({
    el: '#app',
    data: function () {
        return { 
            visible: false,
            dialogVisible:true,
            info:{
                name:'',
                score:0,
                rank:0
            },
            isMeet:false
        }
    },
    computed:{
        clientHeight:function(){
            return document.documentElement.clientWidth
        },
        gateHeight:function(){
            var self = document.getElementById("gate");
            return self.getBoundingClientRect().top + document.documentElement.scrollTop
        },
        geteLeft:function(){
            var self = document.getElementById("gate");
            return self.getBoundingClientRect().left
        }
    },
    methods: {
        getRank:function(){
            rank({name: this.info.name}).then(Response => {
                this.info.rank = Response.data.rank
            }).catch(e => {
                this.$message({
                    message: "获取排名失败" + e,
                    type:"danger"
                })
            })
        },
        submitInfo:function(){
            this.dialogVisible = false
            login({name: this.info.name}).then(Response => {
                this.info.name = Response.data.to
                this.$message({
                    message:"欢迎" + Response.msg
                })
                this.getRank()
            }).catch(e => {
                this.$message({
                    message: "登陆失败" + e,
                    type:"danger"
                })
            })
        },
        moveBall:function(){
            var elem = document.getElementsByClassName("ball")
            var gate = document.getElementById("gate")
            var id = setInterval(frame, 20)
            var tempHeight = 0
            var that = this
            function frame() {
                var space = elem[0].getBoundingClientRect().left - gate.getBoundingClientRect().left
                var hight = elem[0].getBoundingClientRect().top - gate.getBoundingClientRect().top
                if (!that.isMeet && (space <60 && space > -20 && hight < 20 && hight > -20)) {
                    that.isMeet = true
                    score(JSON.stringify({name:that.info.name, score:1})).then(Response => {
                        that.info.score = Response.data.score
                        that.$message({
                            message:"进球成功, 总分数为:" + Response.data.score,
                            type:"success"
                        })
                        that.getRank()
                    }).catch( e => {
                        that.$message({
                            message: "分数统计失败" + e,
                            type:"danger"
                        })
                    })
                }
                if (tempHeight >= that.clientHeight) {
                    elem[0].style.bottom = 0 + ''
                    clearInterval(id)
                } else{
                    tempHeight += 20
                    elem[0].style.bottom = tempHeight + ''
                }
            }
        },
        move: function () {
            var elem = document.getElementById("myBar")
            var width = 0
            var id = setInterval(frame, 20)
            function frame() {
                if (width >= 100) {
                    clearInterval(id)
                    elem.remove()
                } else {
                    width++
                    elem.style.width = width + '%'
                    elem.innerHTML = width * 1 + '%'
                }
            }
            
        },
        moveGate:function(){
            var elem = document.getElementById("gate")
            var id = setInterval(frame, 20)
            var model = false
            var tempValue = 0
            var increase = 1
            function frame() {
                if (tempValue > 100 || tempValue < -100) {
                    increase = -increase
                }
                tempValue += increase
                elem.style.marginLeft = tempValue + 'px'
            }
        },
        clickPeople(){
            this.isMeet = false
            this.moveBall()
        }
    },
    mounted(){
        this.move()
        this.moveGate()
    },
    created() {
        
    }
})