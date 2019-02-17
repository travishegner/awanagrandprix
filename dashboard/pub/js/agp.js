$(function() {
	loadSeasons();
});

function loadSeason(season) {
	$.ajax('/parts/season.html').done(function(tpl){
		$('#content').html("").html(tpl);

		var li = $('ol.breadcrumb li:first').clone();
		li.children('a').html(season.name);

		$('ol.breadcrumb li:first').removeClass('active');
		$('ol.breadcrumb').append(li);
	});
}

function loadSeasons() {
	$.ajax('/parts/seasons.html').done(function(tpl){
		$('#content').html("").html(tpl);

		$('#addseason').click(function() {
			var aj = $.ajax({
				type: "POST",
				url: '/api/seasons',
				data: "season="+$('#season').val(),
				dataType: "json"
			});
			aj.done(function(j) {
				loadSeasonList(); 
			});
			aj.fail(function() {
				console.log("Failed to add season. Does one with that name already exist?");
			});
		});

		$('#season').keyup(function(e) {
			if(e.which == 13) {
				$('#addseason').trigger("click");
			}
		});

		$('#seasons').on('click', 'a', function(e) {
			var t = $(e.currentTarget);
			loadSeason({id: t.attr('data-id'), name: t.html()});
		});

		loadSeasonList();
	});
}

function loadSeasonList() {
	$.ajax('/api/seasons').done(function(seasons) {
		$('#seasons li').not(".template").remove();
		$('#season').val("");

		$.each(seasons, function(i, s) {
			var li = $('#seasons li.template').clone().removeClass('template');
			li.children('a').html(s.name).attr("data-id", s.id);
			$('#seasons').append(li);
		});
	});
}
