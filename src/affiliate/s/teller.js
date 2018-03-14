function validateJoinForm(form){
	var addrInput = form.elements['address'];
	$(addrInput).val($.trim($(addrInput).val()));
	if(!$(addrInput).val()){
		alert('Please input wallet address');
		return false;
	}
	return true
}
function generate(){
	var form = document.forms['joinForm'];
	if(!validateJoinForm(form)){
		return;
	}
	$.ajax({ url: './generate/', method: 'POST', data: $(form).serialize(), dataType: 'json', success: function(obj){
		if(obj){
			if(obj.code===0){
				$('#buyUrl').removeAttr("readonly");
				$('#buyUrl').val(obj.data.buyUrl);
				$('#buyUrlQr').html('<img height="100%" src="/qr-code/?v=2&content='+encodeURIComponent(obj.data.buyUrl)+'"/>')
				$('#buyUrl').attr("readonly","readonly");
                showLayer('generateTrackingUrl');
			}else{
				alert(obj.errmsg);
			}
		}else{
			alert("error"); 		
		}
		}});
}
function myInvitation(){
	var form = document.forms['joinForm'];
	if(!validateJoinForm(form)){
		return;
	}
	$.ajax({ url: '/my-invitation/', method: 'POST', data: $(form).serialize(), dataType: 'json', success: function(obj){
		if(obj){
			if(obj.code===0){
				$.views.settings.delimiters("[[", "]]");
				var tmpl = $.templates('#tmpl')
				$("#myInvitationDiv").html(tmpl.render(obj.data));
				showLayer('myInvitation');
			}else{
				alert(obj.errmsg);
			}
		}else{
			alert("error"); 		
		}
		}});
}

function validateBuyForm(form){
	var addrInput = form.elements['address'];
	$(addrInput).val($.trim($(addrInput).val()));
	if(!$(addrInput).val()){
		alert('Please input wallet address');
		return false;
	}
    var radios = form.elements['currencyType'];
    var valid = false;
    for (var i = 0; i < radios.length; i++){
     if (radios[i].checked){
      valid = true;
      break;
     }
    }
    if(!valid){
    	alert('Please choose cryptocurrency type');
    	return false;
    }
	return true
}
function getAddress(){
    var form = document.forms['buyForm'];
    if(!validateBuyForm(form)){
    	return;
    }
    $.ajax({ url: '/get-address/', method: 'POST', data: $(form).serialize(), dataType: 'json', success: function(obj){
    	if(obj){
    		if(obj.code===0){
				$('#depositAddr').removeAttr("readonly");
				$('#depositAddr').val(obj.data.depositAddr);
				//$('#depositAddrInfo').html(obj.data.first?'Your wallet address is the first time getting address.':'Your wallet address have already got the address.')
				$('#depositAddrQr').html('<img height="100%" src="/qr-code/?v=2&content='+encodeURIComponent(obj.data.depositAddr)+'"/>')
                $('#depositAddr').attr("readonly","readonly");
                $('#depositAddrQr').show();
    			$('#statusRes').hide();
                showLayer('depositAddress');
    		}else{
    			alert(obj.errmsg);
    		}
    	}else{
    		alert("error"); 		
    	}
      }});
}
function checkStatus(){
    var form = document.forms['buyForm'];
    if(!validateBuyForm(form)){
    	return;
    }
    $.ajax({ url: '/check-status/', method: 'POST', data: $(form).serialize(), dataType: 'json', success: function(obj){
    	if(obj){
    		if(obj.code===0){
                $('#depositAddrQr').hide();
    			$('#statusRes').text(obj.data);
    			$('#statusRes').show();
    		}else{
    			alert(obj.errmsg);
    		}
    	}else{
    		alert("error"); 		
    	}
      }});
}

window.setInterval(function(){
	$.ajax({ url: '/get-rate/', method: 'POST',  dataType: 'json', success: function(obj){
    	if(obj){
    		if(obj.code===0){
    			for(var i=0;i<obj.data.length;i++){
    				var node = $('#rate-'+obj.data[i].code);
    				if(!obj.data[i].enabled){
    					var labelNode = $('#li-'+obj.data[i].code);
    					if(labelNode){
    						labelNode.hide();
    					}
    				}else if(node&&obj.data[i].reverse_rate!=node.text()){
						node.text(obj.data[i].reverse_rate);  					
    				}
    			}
    		}else{
    			console.log(obj.errmsg);
    		}
    	}else{
    		console.log("error"); 		
    	}
      }});
},60000)

window.alert = function(message){
    $('#commonAlertMsg').html(message);
    showLayer('commonAlert');
}

function copyToClipboard(id, resultId) {
    var text = $('#'+id).val();
    if(!text){
        return;
    }
    var copyBox = document.createElement('textarea');
    copyBox.style.position = 'fixed';
    copyBox.style.left = '0';
    copyBox.style.top = '0';
    copyBox.style.opacity = '0';
    copyBox.value = text;
    document.body.appendChild(copyBox);
    copyBox.focus();
    copyBox.select();
    document.execCommand('copy');
    document.body.removeChild(copyBox);
    if(resultId){
        $('#'+resultId).html('Copied!');
    }
    setTimeout(function () {
        $('#'+resultId).html('');
      }, 3000);
}