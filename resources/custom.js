$(function () {
    $('div.checkboxes').hide();
});


$( "#filter" ).submit(function( event ) {
    alert( "Handler for .submit() called." );
    event.preventDefault();
  });


$(document).ready(function () {
    $("p").click(function () {
        if ($('div.checkboxes').is(':visible')) {
            $("div.checkboxes").slideToggle();
            $("#filter").submit();
        } else {
            $("div.checkboxes").slideToggle();
        }
    });
});


