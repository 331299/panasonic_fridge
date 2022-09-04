Java.perform(function () {
  var myLog = Java.use("cn.pana.caapp.cmn.communication.MyLog");

  myLog.d.overload("java.lang.String", "java.lang.String").implementation =
    function (strA, strB) {
      console.log(strA, strB);
    };

  myLog.e.overload("java.lang.String", "java.lang.String").implementation =
    function (strA, strB) {
      console.log(strA, strB);
    };

  myLog.i.overload("java.lang.String", "java.lang.String").implementation =
    function (strA, strB) {
      console.log(strA, strB);
    };

  myLog.v.overload("java.lang.String", "java.lang.String").implementation =
    function (strA, strB) {
      console.log(strA, strB);
    };

  myLog.w.overload("java.lang.String", "java.lang.String").implementation =
    function (strA, strB) {
      console.log(strA, strB);
    };
});
