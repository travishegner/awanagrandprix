$(function() {
	$.ajax('/api/seasons').done(function(data){
		$.each(data, function(i, d) {
			var li = $('#seasons li.template').clone().removeClass('template');
			console.log(d)
			li.children('a').html(d.name)
			$('#seasons').append(li)
		});
	});
});
