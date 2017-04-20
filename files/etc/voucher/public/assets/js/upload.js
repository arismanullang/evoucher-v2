var data = new FormData();
jQuery.each(jQuery('#img_url')[0].files, function(i, file) {
    data.append('file-'+i, file);
});

$.ajax({
	url: '/v1/file/upload?token='+token,
	data: data,
	cache: false,
	contentType: false,
	processData: false,
	type: 'POST',
	success: function(data){
	alert(data);
}
});
