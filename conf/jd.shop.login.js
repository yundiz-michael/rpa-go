/*auth:@admin
  name:@jd.shop.login
  desc:@京麦商家登录
 */
function run(page) {
    var tabs = page.selects("div.tabs-header-item");
    tabs[1].click();
    page.waitVisible("#loginFrame");
    var frame = page.select("#loginFrame");
    frame.select("input#loginname").setValue("tianjun-jiangkuo");
    frame.select("input#nloginpwd").setValue("A12345678");
    frame.select("div.login-btn2013").click();
    /*处理登录图片*/
    var qryBigImage = "div.JDJRV-bigimg img";
    var qrySmallImage = "div.JDJRV-smallimg img";
    bigImage = frame.imageReady(qryBigImage);
    smallImage = frame.imageReady(qrySmallImage);
    attrBig = frame.select(qryBigImage).attrsBy(["width", "height"]);
    var count = 3;
    var loginSuccess = false;
    while (count >= 0) {
        var resp = page.client.decodeJDLogin(smallImage.substring(22),
            bigImage.substring(22));
        if (resp.isSuccess) {
            var width = parseFloat(attrBig["width"]);
            var offsetX = (width / parseFloat(resp["w"])) * parseFloat(resp["x"]) - 2;
            frame.mouseDrag("div.JDJRV-slide-btn", offsetX);
            resp = page.waitVisible("span.shop-pageframe-navigation__shop-name", 5);
        }
        loginSuccess = resp.isSuccess;
        if (loginSuccess) break;
        var btnRefreshImage = frame.select("div.JDJRV-img-refresh")
        btnRefreshImage.click();
        bigImage = frame.imageChanged(qryBigImage, bigImage);
        smallImage = frame.imageChanged(qrySmallImage, smallImage);
        page.wait(1000);
        count--;
    }
    return loginSuccess;
}

exports.run = run;



