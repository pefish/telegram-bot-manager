var commands = {
    "/test": {
        desc: "测试命令",
        func: function (args) {
            // console.log(args)
            return "test: " + JSON.stringify(args)
        }
    },
    "/haha": {
        desc: "有点意思",
        func: function (args) {
            return "xixi"
        }
    },
}



