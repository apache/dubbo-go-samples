new Vue({
    el: '#app',
    data: function () {
        return { 
            visible: false,
            rootDom:Object, 
            dialogVisible:true,
            info:{
                name:'',
                uid:'',
                score:0
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
        submitInfo:function(){
            this.dialogVisible = false
            login({name: this.info.name}).then(Response => {
                this.info.uid = Response.to
                console.log(this.info)
                this.$message({
                    message:"欢迎" + Response.message
                })
            }).catch((e) => {});
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
                    score().then(Response => {
                        that.$message({
                            message:"进球成功",
                            type:"success"
                        })
                    })
                }
                if (tempHeight >= that.clientHeight) {
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
        this.rootDom = document.getElementById("app")
    },
    created() {
        
    }
})