$.extend({
showModal: function(id) {
    $(id).modal("show");
},
showTab: function(id) {
    $(id).click();
},
showPreview: function() {
    $.showModal(cS.previewModal);
    title=$(cS.title).val();
    content=$(cS.content).val();
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
                fileNum = (i+1+(pid-1)*cS.numListFiles);
                filesInfo.push('<tr class="file" fileid="'+i+'"><td class="num">'+fileNum+'</td><td class="name">'+$.getFileIcon(file.Type)+' <a href="/file?key='+file.ID+'" target="blank">'+file.Name+'</a></td><td class="Operations"><a fileid="'+i+'" onclick="$.initEdit('+i+');" class="btn btn-mini edit"><i class="icon-pencil"></i> </a><a href="javascript:void(0)" onclick="$.deleteFile('+i+','+pid+');" class="btn btn-mini"><i class="icon-trash"></i> </a></td></tr>');
            });
            $(cS.listData).html(filesInfo.join(''));

            // Insert nav data
            var navInfo = [];
            $(cS.listNav).html("");
            if (page.Nav.ShowPrev) {
                navInfo.push('<li class="previous"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.PrevPageID+');">&larr; Older</a></li>');
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
                navInfo.push('<li class="next"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.NextPageID+');">Newer &rarr;</a></li>');
            }
            $(cS.listNav).html(navInfo.join(''));

        } else {
            $(cS.listData).html("");
            $(cS.listNav).html("");
        }
        $.showMessage(message.Info,"list");
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
initEdit: function(id) {
    fileEdit = files[id];

    var type = fileEdit.Type;
    var src = cS.urlFile+fileEdit.ID;

    if (cS.style == "tab") {
        $(cS.uploadSave).html("Upload").attr("id","edit-save");
        $.showTab(cS.togEdit);
    } else {
        $("#file-upload-modal").modal('show');
    }

    $(cS.editPreview).html($.getPreview(type, src));
    $(cS.editDes).val(fileEdit.Description);
    $(cS.editName).val(fileEdit.Name);
    //$(cS.message).html("");
    $.showMessage("", "edit");

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
deleteFile: function(id, pid){
    fileDelete = files[id];
    //$(cS.message).html("delete file...");
    $.showMessage("delete file...", "list");
    $.post(cS.urlDelete+fileDelete.ID, function(message){
        $.initList(pid);
        $(cS.message).html(message.Info);
    });
},
initUpload: function(){
    if (filesUpload.length>0){
        if (cS.style == "tab") {
            $(cS.uploadSave).html("Upload").attr("id","upload-save");
            $.showTab(cS.togUpload);
        } else {
            $("#file-upload-modal").modal('show');
        }
        
        var filesInfo = [];

        for (var i = 0; i < filesUpload.length; i++) {
            filesInfo.push('<tr><td>'+(i+1)+'</td><td>'+filesUpload[i].name+'</td><td><a href="javascript:void(0)" onclick="$.deleteUploadFile('+i+');" class="btn btn-mini"><i class="icon-trash"></i> </a></td></tr>');
        };

        $(cS.uploadData).html(filesInfo.join(''));
    } else {
        $(cS.uploadData).html("");
        if (cS.style == "tab") {
            $.showTab(cS.togList);
        } else {
            $("#file-upload-modal").modal('hide');
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
    $.showMessage("Uploading...", "upload");
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
            if (cS.style == "tab") {
                $(cS.uploadData).html("");
                $.showTab(cS.togList);
            } else {
                $("#file-upload-modal").modal('hide');
            }

        }
        $.showMessage(message.Info, "upload");
    });
},
showMessage: function(message, q){
    if (cS.style == "tab") {
        id = cS.message;
    } else if (q == "list") {
        id = "#file-upload-modal .message";
    } else {
        id = "#file-upload-modal .message";
    }
    $(id).html(message);
},
saveEditFile: function(){
    var name = $(cS.editName).val();
    var description = $(cS.editDes).val();
    if (name != "" && (name != fileEdit.Name || description != fileEdit.Description)){
        $.post(cS.urlEdit+fileEdit.ID,{name:name,description:description},function(message){
            $(".message-info").html(message.Info);
            $.showMessage(message.Info, "edit")
            $("#file-edit-modal").modal("hide");
            initFileList(1);
        });
        $("#file-edit-modal .message").html("Save changes...");
    } else {
        $("#file-edit-modal .message").html("Make some changes!");
    }
}


});

$(document).ready(function(){
    // preview-modal
    $(cS.preview).click(function(){
        $.showPreview();
    });

    $(cS.previewSave).click(function(){
        $(cS.submit).click()
    });

    // upload-modal
    $(cS.upload).click(function(){
        $.initList(1);
        $.showModal(cS.uploadModal);
    });

    $("#"+cS.file).change(function(){
        filesUpload = document.getElementById(cS.file).files;
        $.initUpload();
    });

    $(cS.uploadModal+" .save").click(function(){
        if (($(this).attr("id"))=="upload-save"){
            $.uploadFiles();
        }
        switch ($(this).attr("id")) {
            case "upload-save":
                $.uploadFiles();
                break
            case "edit-save":
                $.saveEditFile();
                break
        }

    });

});