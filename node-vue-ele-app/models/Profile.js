const mongoose = require("mongoose");
const Schema = mongoose.Schema;

// Create Schema
const ProfileSchema = new Schema({
    type: {
        type: String
    },
    describe: {
        type: String
    },
    income: {
        type: String,
        required: true
    },
    expend: {
        type: String,
        required: true
    },
    cash: {
        type: String,
        required: true
    },
    remark: {
        type: String
    },
    date: {
        type: String,
        default: format(new Date)
    }
})


 function format(Date){
	var Y = Date.getFullYear();
	var M = Date.getMonth() + 1;
		M = M < 10 ? '0' + M : M;// 不够两位补充0
	var D = Date.getDate();
		D = D < 10 ? '0' + D : D;
	var H = Date.getHours();
		H = H < 10 ? '0' + H : H;
	var Mi = Date.getMinutes();
		Mi = Mi < 10 ? '0' + Mi : Mi;
	var S = Date.getSeconds();
		S = S < 10 ? '0' + S : S;
        return Y + '-' + M + '-' + D + ' ' + H + ':' + Mi + ':' + S;
}

module.exports = Profile = mongoose.model("profile", ProfileSchema);