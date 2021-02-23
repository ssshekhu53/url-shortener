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
        let fd=new FormData(form);
        $.ajax({
           type: 'POST',
           url: '/tinyfy',
           data: fd,
           contentType: false,
           processData: false,
           dataType: 'json',
           beforeSend: () => {
               $('.alert').fadeOut();
               $(form).find('button:not(#copy-btn)').attr('disabled', true);
               $(form).find('button:not(#copy-btn)').text('Tinyfying...');
               $('.alert input[name="short_url"]').val('');
               $('.alert').hide();
               $('#copy-btn').html('<i class="far fa-copy"></i>Copy');
               $('#copy-btn').removeClass('btn-success');
               $('#copy-btn').addClass('btn-danger');
               $('#copy-btn').hide();
           },
           success: function (response) {
               console.log(response);
               $('.alert input[name="short-url"]').attr('value', response.short_url);
               $('.alert').fadeIn();
               $(form).find('button:not(#copy-btn)').text('Make TinyURL!');
               $(form).find('button:not(#copy-btn)').attr('disabled', false);
               $('#copy-btn').show();
           }
        });
    }
});

$(document).on('click', '#copy-btn', function(){
    $(`input[name="short-url"]`).select();
    document.execCommand("copy");
    $(this).html('<i class="fas fa-check"></i>Copied');
    $(this).removeClass('btn-danger');
    $(this).addClass('btn-success');
});