/*auth:@admin  name:@jd.shop.day desc:@抓取京东店铺日销售
订单数据格式
{
  "shopName":"店铺名称",
  "logisticsId":"物流编号",
  "logisticsName":"物流名称",
  "customer":{
    "name":"收货人姓名",
    "mobile":"收货人电话",
    "dongdong":"东东号",
  }
  "createdTime":"下单时间",
  "paidTime":"付款时间",
  "amount":"商品总额 float",
  "freightAmount":"运费金额 float",
  "promotionPrice":"促销价格 float",
  "coupon":"优惠券 float",
  "bean":"京豆 float",
  "paidAmount":"应支付金额 float",
  "items":[  //sku列表
     {
     	"id":"货品代码",
     "name":"货品名称",
     "qty":"数量"
     "price":"价格",
     "disaccount":"折扣",
     }
  ]
}
 */
var client
if (client==null) client=new WebClient('jd.com');
var page=client.load("home","https://shop.jd.com");

function checkLogin(){
    var _logout=page.select("//a[contains(text(),'退出')]",2);
    var loginSuccess=true;
    if (_logout==null){
        jdLogin= require("/jd.shop.login");
        loginSuccess=jdLogin.run(page);
    }
    console.log(loginSuccess? "登录成功":"登录失败");
    return loginSuccess;
}

function main(){
    if (!checkLogin()) return;
    /*关闭通知对话框 */
    var e=page.select("div.news-announce i.el-dialog__close",1);
    if (e!=null){
        e.click();
    }
    var e=page.select("//div[@class='shop-pageframe-sidebar__fixed-list']//span[contains(text(),'订单查询与跟踪')]");
    e.click();
    /*等待 近三个月订单大于2条，才开始解析 */
    page.waitMoreThan("div.shopweb-table-wrapper",2);
    var shopName=page.select("span.shop-pageframe-navigation__shop-name").text;
    jdShopOrders= require("/jd.shop.orders");
    var todayHasData=jdShopOrders.parseToday(shopName,page);
    if (!todayHasData) return;
    e=page.select("li.ivu-page-next");
    while (e!=null){
        e.click();
        page.waitMoreThan("div.shopweb-table-wrapper",2);
        todayHasData=jdShopOrders.parseToday(shopName,page);
        if (todayHasData) break;
    }
}

main();