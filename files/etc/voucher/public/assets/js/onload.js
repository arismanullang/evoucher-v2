var user = localStorage.getItem("user");
var token = localStorage.getItem(user);

$( window ).ready(function() {
  console.log(token);
  if(token == null){
    window.location = "http://voucher.apps.id:8889/user/login";
  }

  getSession();
});

function getSession() {
    $.ajax({
      url: 'http://voucher.apps.id:8889/v1/api/get/session?token='+token+'&user='+user,
      type: 'get',
      success: function (data) {
        console.log(data.data);
        if(data.data == false){
          window.location = "http://voucher.apps.id:8889/user/login";
        }
      },
      error:function (data) {
        console.log(data.status);
        if(data.status == 500){
          window.location = "http://voucher.apps.id:8889/user/login";
        }
      }
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
