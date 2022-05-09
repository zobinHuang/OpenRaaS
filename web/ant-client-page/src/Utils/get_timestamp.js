function GetTimestamp() {
    let myDate = new Date();
    let str = myDate.toTimeString();
    let timeStr = str.substring(0,8);
    return timeStr
}

export default GetTimestamp