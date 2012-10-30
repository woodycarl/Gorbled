var cS = {
    // tab
    togList: "a[href='#file-list']",
    togEdit: "a[href='#file-edit']",
    togUpload: "a[href='#file-upload']",
    Edit: "#file-edit",
    List: "#file-list",
    Upload: "#file-upload",

    // content
    file: "file",
    listData: "#file-list .data",
    listNav: ".pageination .pager",
    uploadData: "#file-upload .data",
    editPreview: "#file-edit .preview",
    editName: "#file-edit .name",
    editDes: "#file-edit .description",

    message: "#modal-upload .message",
    modalMessage: "",

    uploadSave: "#modal-upload .save",

    // preview
    previewTitle: "#modal-preview .title",
    previewContent: "#modal-preview .modal-body",
    previewSave: "#modal-preview .save",

    // common
    urlGetData: "/admin/file-data?pid=",
    urlFile: "/file?key=",
    urlDelete: "/admin/file-delete?id=",
    urlNew: "/admin/file-new-url?num=",
    urlDecode: "/decodeContent",
    urlEdit: "/admin/file-edit?id=",

}

var files = [];
var filesUpload = [];
var fileEdit;
var uploadUrls =[];

$(document).ready(function(){

    // preview-modal
    $(tS.preview).click(function(){
        $.showPreview();
    });

    $(cS.previewSave).click(function(){
        $(tS.submit).click()
    });

    // upload-modal
    $(tS.upload).click(function(){
        $.initList(1);
        $.showModal(tS.uploadModal);
    });

    // #file
    $("#"+cS.file).change(function(){
        filesUpload = document.getElementById(cS.file).files;
        $.initUpload();
    });

    $(tS.uploadModal+" .save").click(function(){
        switch ($(this).attr("id")) {
            case "upload-save":
                $.uploadFiles();
                break
            case "edit-save":
                $.saveEditFile();
                break
            case "list-refresh":
                $.initList(1);
                break
        }
    });

    $('a[data-toggle="tab"][href="#file-list"]').on('shown', function (e) {
        $(cS.uploadSave).html(tS.lRefresh).attr("id", "list-refresh");
    });
    $('a[data-toggle="tab"][href="#file-edit"]').on('shown', function (e) {
        $(cS.uploadSave).html(tS.lSaveEdit).attr("id", "edit-save");
    });
    $('a[data-toggle="tab"][href="#file-upload"]').on('shown', function (e) {
        $(cS.uploadSave).html(tS.lUpload).attr("id", "upload-save");
    });

    $(cS.Edit+" .cancel").live("click",function(){
        $(cS.editPreview).html("");
        $(cS.Edit+" .message").html("");
    });

});

$.extend({
showModal: function(id) {
    $(id).modal("show");
},
hideModal: function(id) {
    $(id).modal("hide");
},
showTab: function(id) {
    $(id).click();
},
showPreview: function() {
    $.showModal(tS.previewModal);
    title=$(tS.title).val();
    content=$(tS.content).val();
    $.post(cS.urlDecode,{content:content},function(result){
        $(cS.previewTitle).html(title);
        $(cS.previewContent).html(result);
    });
},
initList: function(pid) {
    uploadUrls = [];
    filesUpload = [];
    files = [];

    $.getJSON(cS.urlGetData+pid, function(message) {
        if (message.Success) {
            var page = jQuery.parseJSON(message.Data);
            files = page.Files;
            
            // Insert file data to table
            var filesInfo = [];
            $.each(page.Files, function(i, file) {
                fileNum = (i+1+(pid-1)*tS.numListFiles);
                var s = '';
                var e = '';
                var d = '';
                var b = 'btn-mini';
                if (tS.style!="tab") {
                    s = '<td class="date">'+file.Date+'</td>';
                    e = tS.lEdit;
                    d = tS.lDelete;
                    b = '';
                }
                filesInfo.push('<tr class="file" fileid="'+i+'"><td class="num">'+fileNum+'</td><td class="name">'+$.getFileIcon(file.Type)+' <a href="/file?key='+file.ID+'" target="blank">'+file.Name+'</a></td>'+s+'<td class="Operations"><a fileid="'+i+'" onclick="$.initEdit('+i+');" class="btn '+b+' edit"><i class="icon-pencil"></i> '+e+'</a><a href="javascript:void(0)" onclick="$.deleteFile('+i+','+pid+');" class="btn '+b+'"><i class="icon-trash"></i> '+d+'</a></td></tr>');
            });
            $(cS.listData).html(filesInfo.join(''));

            // Insert nav data
            var navInfo = [];
            $(cS.listNav).html("");
            if (page.Nav.ShowPrev) {
                navInfo.push('<li class="previous"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.PrevPageID+');">&larr; '+tS.lOlder+'</a></li>');
            }
            if (page.Nav.ShowIDs) {

                $.each(page.Nav.PageIDs, function(i, pageID) {
                    if (pageID.Current) {
                        navInfo.push('<li class="current"><a>'+pageID.Id+'</a></li>');
                    } else {
                        navInfo.push('<li class=""><a href="javascript:void(0)" onclick="$.initList('+pageID.Id+');">'+pageID.Id+'</a></li>');
                    }
                });
            }
            if (page.Nav.ShowNext) {
                navInfo.push('<li class="next"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.NextPageID+');">'+tS.lNewer+' &rarr;</a></li>');
            }
            $(cS.listNav).html(navInfo.join(''));

        } else {
            $(cS.listData).html("");
            $(cS.listNav).html("");
        }
        $.showMessage(message.Info,"list");
    });
},
deleteFile: function(id, pid){
    fileDelete = files[id];
    $.showMessage(tS.lDeleteFile, "list");
    $.post(cS.urlDelete+fileDelete.ID, function(message){
        $.initList(pid);
        $(cS.message).html(message.Info);
    });
},
getFileIcon: function(type) {
    var content;

    if (type.indexOf("image") != -1) {
        content = 'picture';
    } else if (type.indexOf("audio") != -1) {
        content = 'music';
        // headphones | music
    } else if (type.indexOf("video") != -1) {
        content = 'facetime-video'; 
        // film | facetime-video
    } else if (type.indexOf("application/octet-stream") != -1) {
        content = 'gift'; 
        // calendar | th-list | gift
    } else {
        content = 'file';
    }

    return '<i class="icon-'+content+'"></i>';
},
// edit
initEdit: function(id) {
    fileEdit = files[id];
    var type = fileEdit.Type;
    var src = cS.urlFile+fileEdit.ID;

    $(cS.editPreview).html($.getPreview(type, src));
    $(cS.editDes).val(fileEdit.Description);
    $(cS.editName).val(fileEdit.Name);
    $.showMessage("", "edit");

    $.showEdit();
},
showEdit: function(){
    if (tS.style == "tab") {
        $.showTab(cS.togEdit);
    } else {
        $.showModal(cS.Edit);
    }
},
getPreview: function(type, src){
    var preview;
    if (type.indexOf("image") != -1) {
        preview = '<img src="'+src+'" style="max-height:500px;width:auto;" />';
    } else if (type.indexOf("audio") != -1) {
        preview = '<audio src="'+src+'" type="audio/mp3" width="510" height="35" controls="controls" ></audio>';
    } else if (type.indexOf("video") != -1) {
        preview = '<video src="'+src+'" width="510" height="340" controls="controls"></video>';
    } else {
        preview = '<img src="/static/img/readme.png" style="max-height:500px;width:auto;" />';
    }
    return preview;
},
saveEditFile: function(){
    var name = $(cS.editName).val();
    var description = $(cS.editDes).val();
    if (name != "" && (name != fileEdit.Name || description != fileEdit.Description)){
        $.post(cS.urlEdit+fileEdit.ID,{name:name,description:description},function(message){
            $(".message-info").html(message.Info);
            $.showMessage(message.Info, "edit")
            $("#file-edit-modal").modal("hide");
            $.initList(1);
        });
        $("#file-edit-modal .message").html("Save changes...");
    } else {
        $("#file-edit-modal .message").html("Make some changes!");
    }
    $.closeEdit();
},
closeEdit: function(){
    if (tS.style == "tab") {
        $.showTab(cS.togList);
    } else {
        $.hideModal(cS.Edit);
    }
},
initUpload: function(){
    if (filesUpload.length>0){
        // show #file-upload
        if (tS.style == "tab") {
            $.showTab(cS.togUpload);
        } else {
            $(cS.Upload).modal('show');
        }

        var filesInfo = [];

        for (var i = 0; i < filesUpload.length; i++) {
            filesInfo.push('<tr><td>'+(i+1)+'</td><td>'+$.getFileIcon(filesUpload[i].type)+' '+filesUpload[i].name+'</td><td><a href="javascript:void(0)" onclick="$.deleteUploadFile('+i+');" class="btn btn-mini"><i class="icon-trash"></i> </a></td></tr>');
        };

        $(cS.uploadData).html(filesInfo.join(''));
    } else {
        $(cS.uploadData).html("");
        if (tS.style == "tab") {
            $.showTab(cS.togList);
        } else {
            //$(cS.Upload).modal('hide');
        }

    }
},
deleteUploadFile: function(i){
    var tmp = filesUpload;
    filesUpload = [];

    for(var j=0,n=0;j<tmp.length;j++){
        if(j!=i){
            filesUpload[n]=tmp[j];
            n++;
        }
    }
    $.initUpload();
},
uploadFiles: function(){
    $.showMessage(tS.lUploading, "upload");
    $.getJSON(cS.urlNew+filesUpload.length, function(message) {
        // Get uploadUrls
        if (message.Success) {
            uploadUrls = jQuery.parseJSON(message.Data);
            for (var i = 0; i < filesUpload.length; i++) {
                // FormData 对象
                var form = new FormData();
                form.append("file", filesUpload[i]);

                // XMLHttpRequest 对象
                var xhr = new XMLHttpRequest();
                xhr.overrideMimeType("text/html;charset=utf-8");

                if (i<filesUpload.length-1){
                    xhr.open("post", uploadUrls[i], true);
                } else {
                    xhr.open("post", uploadUrls[i], false);
                }

                //xhr.open("post", uploadUrls[i], false);
                xhr.send(form);
            };

            $.initList(1);
            if (tS.style == "tab") {
                $.showTab(cS.togList);
            } else {
                $.hideModal(cS.Upload);
            }
            $(cS.uploadData).html("");

        }
        $.showMessage(message.Info, "upload");
    });
},
showMessage: function(message, q){
    if (tS.style == "tab") {
        id = cS.message;
    } else if (q == "list") {
        id = tS.messageList;
    } else {
        id = "#file-upload-modal .message";
    }
    $(id).html(message);
}

});
