$( document ).ready(function() {

});

function getProfile(key){
	$.ajax({
		url: '/v1/redeem/profile?key='+key,
		type: 'get',
		success: function (data) {
			console.log(data);
		}
	});
}

(function() {
	'use strict';

	$(formAdvanced);

	function formAdvanced() {
		$('.select2').select2();
	}
})();
