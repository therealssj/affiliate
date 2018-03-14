var __welcome = {
    slidesImages: function () {
        var $bigImg = $('#bigImg'),
            $prev = $('._btnPrev'),
            $next = $('._btnNext'),
            $thumbnail = $('._thumbnailBox li'),
            iNow = $thumbnail.length - 1,
            current = 0;

        function showImg(i) {
            var strSrc = $thumbnail.find('img').eq(i).attr('data-src');
            $bigImg.attr('src', strSrc);
            $prev.show();
            $next.show();
        }

        showImg(0);
        $prev
            .on('click', function () {
                var $this = $(this);

                if (current <= 0) {
                    $this.hide();
                    return;
                }
                current--;
                showImg(current);
            });

        $next
            .on('click', function () {
                var $this = $(this);

                if (current >= iNow) {
                    current = iNow;
                    $this.hide();
                    return;
                }
                ++current;
                showImg(current);
            });

        $thumbnail
            .on('click', function () {
                var $this = $(this);

                showImg($this.index());
            });

    }
};

$(function () {
    __welcome.slidesImages();
});

function hashChange(){
    var hash = location.hash;
    if(hash=='#wallet'||hash=='#newsletter'){
        showLayer(hash.substring(1)+'Layer');
    }
}


function signUpNewsletter(){
    var form = document.forms['signUpNewsletterForm'];
    var emailInput = form.elements['email'];
	$(emailInput).val($.trim($(emailInput).val()));
	if(!$(emailInput).val()){
        $('#signUpNewsletterErr').html('Please input your email.');
        $('#signUpNewsletterDiv').addClass('form-group-error');
		return;
    }
    var re = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;
    if(!re.test(String($(emailInput).val()).toLowerCase())){
        $('#signUpNewsletterErr').html('This email is not invalid.');
        $('#signUpNewsletterDiv').addClass('form-group-error');
		return;
    }
    $('#signUpNewsletterDiv').removeClass('form-group-error');
    $.ajax({ url: '/record-newsletter-email/', method: 'POST', data: $(form).serialize(), dataType: 'json', success: function(obj){
    	if(obj){
    		if(obj.code===0){
				alert('Sign-Up Newsletter success.');
    		}else{
                $('#signUpNewsletterErr').html(obj.errmsg);
                $('#signUpNewsletterDiv').addClass('form-group-error');
    		}
    	}else{
    		alert("error"); 		
    	}
      }});
}