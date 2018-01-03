window.alert = function(message, title) {
    if($("#bootstrap-alert-box-modal").length == 0) {
        $("body").append('<div id="bootstrap-alert-box-modal" class="modal fade">\
            <div class="modal-dialog">\
                <div class="modal-content">\
                    <div class="modal-header" style="min-height:40px;">\
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>\
                        <h4 class="modal-title"></h4>\
                    </div>\
                    <div class="modal-body"><p></p></div>\
                    <div class="modal-footer">\
                        <a href="#" data-dismiss="modal" class="btn btn-default">Close</a>\
                    </div>\
                </div>\
            </div>\
        </div>');
    }
    $("#bootstrap-alert-box-modal .modal-header h4").text(title || "Alert");
    $("#bootstrap-alert-box-modal .modal-body p").text(message || "");
    $("#bootstrap-alert-box-modal").modal('show');
};

function showTooltip(elem,msg){
	$(elem).tooltip('hide')
    .attr('data-original-title', msg)
    .tooltip('fixTitle')
    .tooltip('show');
}
function fallbackMessage(action){
	var actionMsg='';
	var actionKey=(action==='cut'?'X':'C');
	if(/iPhone|iPad/i.test(navigator.userAgent)){
		actionMsg='No support :(';
	}else if(/Mac/i.test(navigator.userAgent)){
		actionMsg='Press ⌘-'+ actionKey+' to '+ action;
	}else{
		actionMsg='Press Ctrl-'+ actionKey+' to '+ action;
	}
	return actionMsg;
}