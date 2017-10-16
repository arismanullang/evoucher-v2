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

function uploadTest() {
	$('.sendButton').attr("disabled","disabled");
	id = findGetParameter("id");
	var headerImage = null;
	var voucherImage = null;
	var footerImage = null;

	var headerImageUrl = "";
	var voucherImageUrl = "";
	var footerImageUrl = "";

	files.forEach(function (file) {
		switch(file.key){
			case "dropzone-header":
				headerImage = file.value;
				break;
			case "dropzone-voucher":
				voucherImage = file.value;
				break;
			case "dropzone-footer":
				footerImage = file.value;
				break;
		}
	});

	var n = "header";
	var interval = setInterval(function () {
		switch (n){
			case "header":
				console.log(headerImage);
				$('#modal-loader').modal('hide');
				if (headerImage != null) {
					var formData = new FormData();
					formData.append('image-url', headerImage);

					$('#image-header-modal').modal({backdrop: 'static',keyboard: false}, 'show');

					jQuery.ajax({
						url: '/file/upload',
						type: "POST",
						processData: false,
						contentType: false,
						data: formData,
						success: function (data) {
							headerImageUrl = data.data;
						}
					});
				}
				n = "voucher";
				break;
			case "voucher":
				console.log(voucherImage);

				$('#image-header-modal').modal('hide');
				if (voucherImage != null) {
					var formData = new FormData();
					formData.append('image-url', voucherImage);
					$('#image-voucher-modal').modal({backdrop: 'static',keyboard: false}, 'show');

					jQuery.ajax({
						url: '/file/upload',
						type: "POST",
						processData: false,
						contentType: false,
						data: formData,
						success: function (data) {
							voucherImageUrl = data.data;
						}
					});
				}
				n = "footer";
				break;
			case "footer":
				console.log(footerImage);

				$('#image-voucher-modal').modal('hide');
				if (footerImage != null) {
					var formData = new FormData();
					$('#image-footer-modal').modal({backdrop: 'static',keyboard: false}, 'show');

					formData.append('image-url', footerImage);

					jQuery.ajax({
						url: '/file/upload',
						type: "POST",
						processData: false,
						contentType: false,
						data: formData,
						success: function (data) {
							footerImageUrl = data.data;
						}
					});
				}
				n = "lalala";
				break;
			default:
				console.log(headerImageUrl);
				console.log(voucherImageUrl);
				console.log(footerImageUrl);
				if( (headerImage == null && headerImageUrl == "") || (voucherImage != null && voucherImageUrl == "") || (footerImage != null && footerImageUrl == "") ){
					break;
				}

				var campaign = {
					program_id: id,
					image_header: headerImageUrl,
					image_voucher: voucherImageUrl,
					image_footer: footerImageUrl,
				};

				$.ajax({
					url: '/v1/ui/campaign/create?token='+token,
					type: 'post',
					dataType: 'json',
					contentType: "application/json",
					data: JSON.stringify(campaign),
					success: function () {
						$('#image-footer-modal').modal('hide');
						swal("Upload Success");
						clearInterval(interval);
						$('.sendButton').removeAttr("disabled");
					},
					error: function (data) {
						var a = JSON.parse(data.responseText);
						swal("Error", a.errors.detail);
						clearInterval(interval);
					}
				});
				break;
		}
	}, 2000);

}

function generateVoucher() {
	var id = $('#program-id').val();
	swal("Sending Voucher");
	$.ajax({
		url: '/v1/ui/voucher/send-voucher?program=' + id + '&token=' + token,
		type: 'post',
		success: function (data) {
			window.location = "/program/search";
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
