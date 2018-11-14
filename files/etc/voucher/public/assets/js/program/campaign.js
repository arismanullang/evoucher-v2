var files = [];
var bool = false;
$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	var id = findGetParameter('id');
	$('#program-id').val(id);
	getProgram(id);

	var allDropzones = document.querySelectorAll('.dropzone');
	for (var i = 0; i < allDropzones.length; i++) {
		var thisDropzone = allDropzones[i];

		new Dropzone(thisDropzone, {
			url: "/file/post",
			autoProcessQueue: false,
			uploadMultiple: false,
			maxFiles: 1,
			dictDefaultMessage: '<em class="ion-upload icon-24 block"></em>', // default messages before first drop
			paramName: 'file', // The name that will be used to transfer the file
			maxFilesize: 2, // MB
			previewTemplate : "<div class=\"dz-preview dz-file-preview\">" +
			"</div>",
			accept: function(file, done) {
				if (file.type != 'image/jpeg') {
					if(this.files.length > 1){
						this.removeAllFiles();
						this.addFile(file);
					}
					done('Must be jpg file.');
				} else {
					if(this.files.length > 1){
						this.removeAllFiles();
						this.addFile(file);
					}
					done();
				}
			},
			init: function() {
				var dzHandler = this;

				this.on('addedfile', function(file) {
					var key = this.element.getAttribute("id");
					files.push({key : key, value : file});

					console.log(file);
					switch(key) {
						case "dropzone-header":
							$("#filename-header").val(file.name);
							break;
						case 'dropzone-voucher':
							$("#filename-voucher").val(file.name);
							break;
						case 'dropzone-background':
							$("#filename-background").val(file.name);
							break;
						case 'dropzone-footer':
							$("#filename-footer").val(file.name);
							break;
						default:
							break;
					}
				});
				this.on('removedfile', function(file) {
					//console.log('Removed file: ' + file.name);
				});
				this.on('sendingmultiple', function() {

				});
				this.on('successmultiple', function(/*files, response*/) {

				});
				this.on('errormultiple', function(/*files, response*/) {
					console.log("error");
				});
				this.on('maxfilesexceeded', function(file) {
					console.log(this.element);
					console.log(this);
					this.removeAllFiles(true);
					this.addFile(file);
				});
			}
		});
	}
});

function preview(){
	readURL(files);

	var content  = $("#content-email").summernote('code');
	console.log("Content : "+ content);
	if(content == "" || content == "<p><br></p>" || content == "<br>"){
		content  = 'Nantikan program-program Digital Voucher menarik lainnya.'
	}
	$("#content").html(content);
	$('#modal-loader').modal( 'show');
}

function readURL(input) {
	input.forEach( function (element) {
		var reader = new FileReader();

		switch(element.key) {
			case "dropzone-header":
				reader.onload = function (e) {
					$('#image-header').attr('src', e.target.result);
				};

				reader.readAsDataURL(element.value);
				break;
			case 'dropzone-voucher':
				reader.onload = function (e) {
					$('#image-voucher').attr('src', e.target.result);
				};

				reader.readAsDataURL(element.value);
				break;
			case 'dropzone-background':
				reader.onload = function (e) {
					$('#image-background').attr('src', e.target.result);
				};

				reader.readAsDataURL(element.value);
				break;
			case 'dropzone-footer':
				reader.onload = function (e) {
					$('#image-footer').attr('src', e.target.result);
				};

				reader.readAsDataURL(element.value);
				break;
			default:
				break;
		}
	});
}

function getProgram(id) {
	$.ajax({
		url: '/v1/ui/program/detail?id=' + id + '&token=' + token,
		type: 'get',
		success: function (data) {
			var result = data.data[0];

			// Program
			$('#program-name').html(result.name);
			$('#modal-program-name').html(result.name);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}


function upload() {
	$('.sendButton').attr("disabled","disabled");
	id = findGetParameter("id");

	var imgObj = new Object();
	imgObj.headerImage = null;
	imgObj.voucherImage = null;
	imgObj.footerImage = null;

	files.forEach(function (file) {
		switch(file.key){
			case "dropzone-header":
				imgObj['headerImage'] = file.value;
				break;
			case "dropzone-voucher":
				imgObj['voucherImage'] = file.value;
				break;
			case "dropzone-footer":
				imgObj['footerImage'] = file.value;
				break;
		}
	});

	uploadHeaderImage(imgObj);

}

function uploadHeaderImage(imgObj){
	var imgUrlObj = new Object();
	imgUrlObj.headerImageUrl = "";
	imgUrlObj.voucherImageUrl = "";
	imgUrlObj.footerImageUrl = "";

	console.log(imgObj.headerImage);
				$('#modal-loader').modal('hide');
				if (imgObj.headerImage != null) {
					var formData = new FormData();
					formData.append('image-url', imgObj.headerImage);

					$('#image-header-modal').modal({backdrop: 'static',keyboard: false}, 'show');

					jQuery.ajax({
						url: '/file/upload',
						type: "POST",
						processData: false,
						contentType: false,
						data: formData,
						success: function (data) {
							imgUrlObj['headerImageUrl'] = data.data;
						},
						complete: function(data){
							uploadVoucherImage(imgObj, imgUrlObj);
						}
					});
				} else {
					uploadVoucherImage(imgObj, imgUrlObj);
				}
}

function uploadVoucherImage(imgObj, imgUrlObj){
	console.log(imgObj.voucherImage);

	$('#image-header-modal').modal('hide');
	if (imgObj.voucherImage != null) {
		var formData = new FormData();
		formData.append('image-url', imgObj.voucherImage);
		$('#image-voucher-modal').modal({backdrop: 'static',keyboard: false}, 'show');

		jQuery.ajax({
			url: '/file/upload',
			type: "POST",
			processData: false,
			contentType: false,
			data: formData,
			success: function (data) {
				imgUrlObj['voucherImageUrl'] = data.data;
			},
			complete: function (data) {
				uploadFooterImage(imgObj, imgUrlObj);
			}
		});
	} else {
		uploadFooterImage(imgObj, imgUrlObj);
	}
}

function uploadFooterImage(imgObj, imgUrlObj){
	console.log(imgObj.footerImage);

	$('#image-voucher-modal').modal('hide');
	if (imgObj.footerImage != null) {
		var formData = new FormData();
		$('#image-footer-modal').modal({backdrop: 'static',keyboard: false}, 'show');

		formData.append('image-url', imgObj.footerImage);

		jQuery.ajax({
			url: '/file/upload',
			type: "POST",
			processData: false,
			contentType: false,
			data: formData,
			success: function (data) {
				imgUrlObj['footerImageUrl'] = data.data;
			},
			complete: function (data) {
				createCampaign(imgUrlObj);
			}
		});
	}else {
		createCampaign(imgUrlObj);
	}
}

function createCampaign(imgUrlObj){
	console.log(imgUrlObj.headerImageUrl);
	console.log(imgUrlObj.voucherImageUrl);
  console.log(imgUrlObj.footerImageUrl);

  var content  = $("#content-email").summernote('code');
	console.log("Content : "+ content);
	if(content == "" || content == "<p><br></p>" || content == "<br>"){
		content  = 'Nantikan program-program Digital Voucher menarik lainnya.'
	}

	var campaign = {
		program_id: id,
		email_subject: $("#subject-email").val(),
    email_sender: $("#sender-email").val(),
		email_content: content,
		image_header: imgUrlObj.headerImageUrl,
		image_voucher: imgUrlObj.voucherImageUrl,
		image_footer: imgUrlObj.footerImageUrl,
	};

	$.ajax({
		url: '/v2/ui/campaign/create?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(campaign),
		success: function () {
			$('#image-footer-modal').modal('hide');
			swal({
					title: 'Success',
					text: 'Campaign Created',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					// closeOnConfirm: false
					closeOnConfirm: true
				},
				function() {
					window.location = "/program/check?id="+id;
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

(function () {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.summernote').each(function () {
			$(this).summernote({
				height: 380,
				placeholder: 'Any Message...',
				callbacks: {
					onPaste: function (e) {
						var bufferText = ((e.originalEvent || e).clipboardData || window.clipboardData).getData('Text');

						e.preventDefault();

						// Firefox fix
						setTimeout(function () {
							document.execCommand('insertText', false, bufferText);
						}, 10);
					}
				}
			});
		});
	}

})();
