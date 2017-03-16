/*!
 *
 * Centric - Bootstrap Admin Template
 *
 * Version: 1.1
 * Author: @themicon_co
 * Website: http://themicon.co
 * License: https://wrapbootstrap.com/help/licenses
 *
 */

// APP START
// -----------------------------------

(function() {
    'use strict';

    // Disable warning "Synchronous XMLHttpRequest on the main thread is deprecated.."
    $.ajaxPrefilter(function(options, originalOptions, jqXHR) {
        options.async = true;
    });

})();