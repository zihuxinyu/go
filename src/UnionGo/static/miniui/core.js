/**
 * Created by weibaohui on 14-5-24.
 */


function GetParams(url, c) {
    if (!url) url = location.href;
    if (!c) c = "?";
    url = url.split(c)[1];
    var params = {};
    if (url) {
        var us = url.split("&");
        for (var i = 0, l = us.length; i < l; i++) {
            var ps = us[i].split("=");
            params[ps[0]] = decodeURIComponent(ps[1]);
        }
    }
    return params;
}


function onSkinChange(skin) {
    //mini.Cookie.set('miniuiSkin', skin);
    mini.Cookie.set('miniuiSkin', skin, 100);//100天过期的话，可以保持皮肤切换
    window.location.reload()
}
function AddCSSLink(id, url, doc) {
    doc = doc || document;
    var link = doc.createElement("link");
    link.id = id;
    link.setAttribute("rel", "stylesheet");
    link.setAttribute("type", "text/css");
    link.setAttribute("href", url);

    var heads = doc.getElementsByTagName("head");
    if (heads.length)
        heads[0].appendChild(link);
    else
        doc.documentElement.appendChild(link);
}


/**
 * 操作成功的提示
 * @param str
 */
function success(str) {
    mini.showMessageBox({
        showModal: true,
        width: 250,
        title: "提示",
        iconCls: "mini-messagebox-info",
        message: str,
        timeout: 1500,
        x: "center",
        y: "middle"
    });
}

/*****************提示方法******************/
/**
 * alert的提示
 * @param str
 */
function alert(str) {
    mini.showMessageBox({
        showModal: true,
        width: 250,
        title: "提示",
        iconCls: "mini-messagebox-warning",
        message: str,
        timeout: 1500,
        x: "center",
        y: "middle"
    });
}
/**
 * 出现错误的提示
 * @param str
 */
function error(str) {

    mini.showMessageBox({
        showModal: true,
        width: 250,
        title: "提示",
        iconCls: "mini-messagebox-error",
        message: str,
        timeout: 1500,
        x: "center",
        y: "middle"
    });

}
/*****************提示方法******************/
/*****************页面方法******************/

/**
 * 页面弹出调用窗口回执方法
 * @param action
 * @returns {*}
 * @constructor
 */
function CloseWindow(action) {
    if (window.CloseOwnerWindow) return window.CloseOwnerWindow(action);
    else window.close();
}

function onOk() {
    CloseWindow("ok");
}
function onCancel() {
    CloseWindow("cancel");
}
/*****************页面方法******************/
/****************renderers*******************/
function onDateRenderer(e) {
    var date = mini.parseDate(e.value)
    var ss = mini.formatDate(date, "yyyy-MM-dd HH:mm:ss")
    return ss;
}
/****************renderers*******************/

/***********grid op start***********/

function addRow(htmlguidid) {
    var grid = mini.get(htmlguidid);
    grid.addRow({}, 0);
    grid.beginEditCell({}, 0);
}
function removeRow(htmlguidid) {
    var grid = mini.get(htmlguidid);
    var rows = grid.getSelecteds();
    if (rows.length > 0) {
        grid.removeRows(rows, true);
    }
}


/**
 * 保存Grid数据
 * @param grid:要操作的grid
 * @param posturl:接收数据的url
 * @returns {boolean}
 */
function saveGrid(gridid, posturl) {
    var grid = mini.get(gridid);
    var data = grid.getChanges(null, true);

    var json = mini.encode(data);
    
    if (data == "") {
        error('没有数据需要保存');
        return false;
    }
    var msgid = mini.loading("数据保存中，请稍后......", "保存数据");
    $.ajax({
        url: posturl,
        data: { data: json },
        type: "post",
        success: function (text) {
            mini.hideMessageBox(msgid);
            success(text);
            grid.reload();
        },
        error: function (jqXHR, textStatus, errorThrown) {
            mini.hideMessageBox(msgid);
            error(jqXHR.responseText);
        }
    });
}
/***********grid op end***********/
