var token = localStorage.getItem("token");
var error = true;
$( window ).ready(function() {
  if(token == null){
    window.location = "/user/login";
  }

  $( 'input[check="true"]' ).focusout(function() {
    if($(this).val() == ""){
      $(this).addClass("error");
      $(this).parent().closest('div').addClass("input-error");
      error = true;
    }
    else{
      $(this).removeClass("error");
      $(this).parent().closest('div').removeClass("input-error");
    }
  });

  $( 'input' ).attr("autocomplete","off");
  $( '#token' ).val(token);

  getSession();
  setSideNavBar();
  setValidationMessage();

  $('[data-toggle="tooltip"]').tooltip();
});

function setSideNavBar() {
	var ui = localStorage.getItem("ui").split(",");
	var li = $("#sidenav").find("a");
	var bool = false;
	for(var i = 0; i < li.length; i++){
		for(var y = 0; y < ui.length-1; y ++) {
			bool = false;

			if (li[i].getAttribute("ui").match(ui[y])) {
				bool = true;
				break;
			}

		}
		if(bool){
			$(li[i]).parent('li').attr("style", "display:block");
		}
	}
}

function getSession() {
	if(window.location.pathname == "program/campaign"){
		$.ajax({
			url: '/v1/ui/token/check?token='+token+'&url='+window.location.pathname,
			type: 'get',
			success: function (data) {
				if(data.data == false){
					logOut();
				}
			},
			error:function (data) {
				logOut();
			}
		});
	}
}

function logOut() {
  localStorage.clear();
  window.location = "/user/login";
}

function addDecimalPoints(value) {
    var input = " "+value;
    var inputValue = input.replace('.', '').split("").reverse().join(""); // reverse
    var newValue = '';
    for (var i = 0; i < inputValue.length; i++) {
        if (i % 3 == 0 && i != 0 && i != inputValue.length-1) {
            newValue += '.';
        }
        newValue += inputValue[i];
    }
    return newValue.split("").reverse().join("");
}

function toTwoDigit(val){
	if (val < 10){
		return '0'+val;
	}
	else {
		return val;
	}
}

function toTitleCase(str)
{
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
    .substr(1)
        .split("&")
        .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    });
    return result;
}

function setValidationMessage() {
	var validatorRequired = "Please fill this field.";
	var validatorRemote= "Please fix this field.";
	var validatorEmail= "Please enter a valid email address.";
	var validatorUrl= "Please enter a valid URL.";
	var validatorDate= "Please enter a valid date.";
	var validatorDateISO= "Please enter a valid date (ISO).";
	var validatorNumber= "Please enter a valid number.";
	var validatorDigits= "Please enter only digits.";
	var validatorCreditcard= "Please enter a valid credit card number.";
	var validatorEqualTo= "Please enter the same value again.";
	var validatorAccept= "Please enter a value with a valid extension.";
	var validatorMaxlength= "Please enter no more than {0} characters.";
	var validatorMinlength= "Please enter at least {0} characters.";
	var validatorRangelength= "Please enter a value between {0} and {1} characters long.";
	var validatorRange= "Please enter a value between {0} and {1}.";
	var validatorMax= "Please enter a value less than or equal to {0}.";
	var validatorMin= "Please enter a value greater than or equal to {0}.";

	jQuery.extend(jQuery.validator.messages, {
		required: validatorRequired,
		remote: validatorRemote,
		email: validatorEmail,
		url: validatorUrl,
		date: validatorDate,
		dateISO: validatorDateISO,
		number: validatorNumber,
		digits: validatorDigits,
		creditcard: validatorCreditcard,
		equalTo: validatorEqualTo,
		accept: validatorAccept,
		maxlength: jQuery.validator.format(validatorMaxlength),
		minlength: jQuery.validator.format(validatorMinlength),
		rangelength: jQuery.validator.format(validatorRangelength),
		range: jQuery.validator.format(validatorRange),
		max: jQuery.validator.format(validatorMax),
		min: jQuery.validator.format(validatorMin)
	});
}

(function() {
    'use strict';

    $(sidebarNav);

    function sidebarNav() {

        var $sidebarNav = $('.sidebar-nav');
        var $sidebarContent = $('.sidebar-content');

        activate($sidebarNav);

        $sidebarNav.on('click', function(event) {
            var item = getItemElement(event);
            // check click is on a tag
            if (!item) return;

            var ele = $(item),
                liparent = ele.parent()[0];

            var lis = ele.parent().parent().children(); // markup: ul > li > a
            // remove .active from childs
            lis.find('li').removeClass('active');
            // remove .active from siblings ()
            $.each(lis, function(idx, li) {
                if (li !== liparent)
                    $(li).removeClass('active');
            });

            var next = ele.next();
            if (next.length && next[0].tagName === 'UL') {
                ele.parent().toggleClass('active');
                event.preventDefault();
            }
        });

        // find the a element in click context
        // doesn't check deeply, asumens two levels only
        function getItemElement(event) {
            var element = event.target,
                parent = element.parentNode;
            if (element.tagName.toLowerCase() === 'a') return element;
            if (parent.tagName.toLowerCase() === 'a') return parent;
            if (parent.parentNode.tagName.toLowerCase() === 'a') return parent.parentNode;
        }

        function activate(sidebar) {
            sidebar.find('a').each(function() {
                var href = $(this).attr('href').replace('#', '');
                if (href !== '' && window.location.href.indexOf(href) >= 0) {
                    var item = $(this).parents('li').addClass('active');
                    // Animate scrolling to focus active item
                    // $sidebarContent.animate({
                    //     scrollTop: $sidebarContent.scrollTop() + item.position().top
                    // }, 1200);
                    return false; // exit foreach
                }
            });
        }

        var layoutContainer = $('.layout-container');
        var $body = $('body');
        // Handler to toggle sidebar visibility on mobile
        $('#sidebar-toggler').click(function(e) {
            e.preventDefault();
            layoutContainer.toggleClass('sidebar-visible');
            // toggle icon state
            $(this).parent().toggleClass('active');
        });
        // Close sidebar when click on backdrop
        $('.sidebar-layout-obfuscator').click(function(e) {
            e.preventDefault();
            layoutContainer.removeClass('sidebar-visible');
            // restore icon
            $('#sidebar-toggler').parent().removeClass('active');
        });

        // Handler to toggle sidebar visibility on desktop
        $('#offcanvas-toggler').click(function(e) {
            e.preventDefault();
            $body.toggleClass('offcanvas-visible');
            // toggle icon state
            $(this).parent().toggleClass('active');
        });

        // remove desktop offcanvas when app changes to mobile
        // so when it returns, the sidebar is shown again
        window.addEventListener('resize', function() {
            if (window.innerWidth < 768) {
                $body.removeClass('offcanvas-visible');
                $('#offcanvas-toggler').parent().addClass('active');
            }
        });

    }

})();

function errorPlacementInput(error, element) {
	if( element.parent().parent().is('.mda-input-group') ) {
		error.insertAfter(element.parent().parent()); // insert at the end of group
		element.focus();
	}
	else if( element.parent().is('.mda-form-control') ) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
		element.focus();
	}
	else if( element.parent().is('.input-group') ) {
		error.insertAfter(element.parent()); // insert after .mda-form-control
		element.focus();
	}
	else if ( element.is(':radio') || element.is(':checkbox')) {
		error.insertAfter(element.parent().parent().parent().parent().parent().find(".control-label"));
		$("input[name=partner]").removeClass('error');
		element.focus();
	}
	else {
		error.insertAfter(element);
		element.focus();
	}
}

(function() {
    'use strict';

    $(initHeader);

    function initHeader() {

        // Search modal
        var modalSearch = $('.modal-search');
        $('#header-search').on('click', function(e) {
            e.preventDefault();
            modalSearch
                .on('show.bs.modal', function() {
                    // Add class for white backdrop
                    $('body').addClass('modal-backdrop-soft');
                })
                .on('hidden.bs.modal', function() {
                    // Remove class for white backdrop (if not will affect future modals)
                    $('body').removeClass('modal-backdrop-soft');
                })
                .on('shown.bs.modal', function() {
                    // Auto focus the search input
                    $('.header-input-search').focus();
                })
                .modal()
                ;
        });

        // Settings modal
        var modalSettings = $('.modal-settings');
        $('#header-settings').on('click', function(){
            modalSettings
                .on('show.bs.modal', function() {
                    // Add class for soft backdrop
                    $('body').addClass('modal-backdrop-soft');
                })
                .on('hidden.bs.modal', function() {
                    // Remove class for soft backdrop (if not will affect future modals)
                    $('body').removeClass('modal-backdrop-soft');
                })
                .modal()
                ;
        });

    }

})();
