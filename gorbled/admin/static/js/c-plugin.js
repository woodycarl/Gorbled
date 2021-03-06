var cS = {
	// tab
	togList: "a[href='#file-list']",
	togEdit: "a[href='#file-edit']",
	togUpload: "a[href='#file-upload']",

	Edit: "#file-edit",
	editPreview: "#file-edit .preview",
	editName: "#file-edit .name",
	editID: "#file-edit .id",
	editDes: "#file-edit .description",
	editMessage: "#file-edit .message",

	List: "#file-list",
	listData: "#file-list .data",
	listNav: ".pagination ul",
	listMessage: "#file-list .message",

	Upload: "#file-upload",
	uploadData: "#file-upload .data",
	uploadSave: "#modal-upload .save",
	uploadMessage: "#file-upload .message",

	// content
	file: "file",

	modalMessage: "#modal-upload .message",
	message: ".manage-bar .manage-info",

	// preview
	previewTitle: "#modal-preview .title",
	previewContent: "#modal-preview .modal-body",
	previewSave: "#modal-preview .save",

	// common
	urlGetData: "/admin/file/data/",
	urlFile: "/file/",
	urlDelete: "/admin/file/delete/",
	urlNew: "/admin/file/new-url/",
	urlDecode: "/decodeContent",
	urlEdit: "/admin/file/edit/",

}

var files = [];
var filesUpload = [];
var fileEdit;
var uploadUrls =[];

$(document).ready(function(){

	// preview-modal
	$(tS.btnPreview).click(function(){
		$.showPreview();
	});

	$(cS.previewSave).click(function(){
		$(tS.submit).click()
	});

	// upload-modal
	$(tS.btnUpload).click(function(){
		$.initList(1);
		$.showModal(tS.uploadModal);
	});

	// #file
	$("#"+cS.file).change(function(){
		filesUpload = document.getElementById(cS.file).files;
		$.initUpload();
	});

	$(cS.uploadSave).click(function(){
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
				filesInfo.push('<tr class="file" fileid="'+i+'"><td class="num">'+fileNum+'</td><td class="name">'+$.getFileIcon(file.Type)+' <a href="'+cS.urlFile+file.ID+'" target="blank">'+file.Name+'</a></td>'+s+'<td class="Operations"><a fileid="'+i+'" onclick="$.initEdit('+i+');" class="btn '+b+' edit"><i class="icon-pencil"></i> '+e+'</a><a  onclick="$.deleteFile('+i+','+pid+');" class="btn '+b+'"><i class="icon-trash"></i> '+d+'</a></td></tr>');
			});
			$(cS.listData).html(filesInfo.join(''));

			// Insert nav data
			var navInfo = [];
			$(cS.listNav).html("");
			if (page.Nav.ShowPrev) {
				//navInfo.push('<li class="previous"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.PrevPageID+');">&larr; '+tS.lOlder+'</a></li>');
				navInfo.push('<li><a href="javascript:void(0)" onclick="$.initList('+page.Nav.PrevPageID+');">«</a></li>');
			}
			$.each(page.Nav.PageIDs, function(i, pageID) {
				if (pageID.Current) {
					//navInfo.push('<li class="current"><a>'+pageID.Id+'</a></li>');
					navInfo.push('<li class="active"><a>'+pageID.Id+'</a></li>');
				} else {
					//navInfo.push('<li class=""><a href="javascript:void(0)" onclick="$.initList('+pageID.Id+');">'+pageID.Id+'</a></li>');
					navInfo.push('<li><a href="javascript:void(0)" onclick="$.initList('+pageID.Id+');">'+pageID.Id+'</a></li>');
				}
			});
			if (page.Nav.ShowNext) {
				//navInfo.push('<li class="next"><a href="javascript:void(0)" onclick="$.initList('+page.Nav.NextPageID+');">'+tS.lNewer+' &rarr;</a></li>');
				navInfo.push('<li><a href="javascript:void(0)" onclick="$.initList('+page.Nav.NextPageID+');">»</a></li>');
			}
			$(cS.listNav).html(navInfo.join(''));

			$.showMessage(message.Info, "list", "success");
		} else {
			$(cS.listData).html("");
			$(cS.listNav).html("");
			$.showMessage(message.Info, "list", "error");
		}
	});
},
deleteFile: function(id, pid){
	fileDelete = files[id];
	$.showMessage(tS.lDeleteFile, "list", "alert");
	$.post(cS.urlDelete+fileDelete.ID, function(message){
		message = jQuery.parseJSON(message);

		$.initList(pid);
		$.showMessage(message.Info, "list", message.Success);
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

	$(cS.editPreview).html($.getPreview(type, src));
	$(cS.editDes).val(fileEdit.Description);
	$(cS.editName).val(fileEdit.Name);
	$(cS.editID).val(fileEdit.ID);
	$.showMessage("", "edit", true);

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
	var id = $(cS.editID).val();
	var name = $(cS.editName).val();
	var description = $(cS.editDes).val();
	if (id != "" && (id != fileEdit.ID || name != fileEdit.Name || description != fileEdit.Description)){
		$.getJSON(cS.urlEdit+fileEdit.ID,{id:id,name:name,description:description},function(message){
			$.showMessage(message.Info, "edit", "success")
			$.initList(1);
			$.closeEdit();
		});
		$.showMessage("Save changes...", "edit", "alert")
	} else {
		$.showMessage("Make some changes!", "edit", "error")
	}
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
	$.showMessage(tS.lDeleteFile, "upload", "alert");
	for(var j=0,n=0;j<tmp.length;j++){
		if(j!=i){
			filesUpload[n]=tmp[j];
			n++;
		}
	}
	$.initUpload();
},
uploadFiles: function(){
	$.showMessage(tS.lUploading, "upload", true);
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

			$.showMessage(message.Info, "upload", "success");
		} else {
			$.showMessage(message.Info, "upload", "alert");
		}
		
	});
},
showMessage: function(message, r, state){
	//var result;
	//var id ="";
	if (tS.style == "tab") {
		if (r=="list" || r=="edit" || r=="upload") {
			id = cS.modalMessage;
		} else {
			id = cS.message;
		}

	} else {
		if (r=="edit") {
			id = cS.editMessage;
		} else if (r=="upload") {
			id = cS.uploadMessage;
		} else {
			id = cS.message;
		}
	}
	if (message=="") {
		result = "";
	} else {
		result = '<div class="alert alert-'+state+'"><a class="close" data-dismiss="alert">×</a>'+message+'</div>';
	}
	$(id).html(result);
	setTimeout("$(id).html('')",3000);
}

});
