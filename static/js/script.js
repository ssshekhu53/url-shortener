$(document).ready(function(){
});

$('#short-form').validate({
    rules: {
        urlField: 'required'
    },
    messages: {
        urlField: "Please enter a URL"
    },
    submitHandler: function (form) {
        $('.alert').fadeOut();
        $(form).find('button').attr('disabled', true);
        $(form).find('button').text('Tinyfying...');
        let fd=new FormData(form);
        $.ajax({
           type: 'POST',
           url: '/tinyfy',
           data: fd,
           contentType: false,
           processData: false,
           dataType: 'json',
           success: function (response) {
               console.log(response);
               $('.alert').text(response.short_url);
               $('.alert').fadeIn();
               $(form).find('button').text('Make TinyURL!');
               $(form).find('button').attr('disabled', false);
           }
        });
    }
});