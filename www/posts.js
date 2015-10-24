
var events = [];
var news = [];
var aircraft = [];
var vehicles = [];
var projects = [];

$(document).ready(function () {
	

	var url = '/get_posts';
	$.getJSON(url, function(posts) {
		console.log(posts);
				
		for (var pst in posts) {
			var post = posts[pst];
			post.n_detail = post.detail;
			post.n_specs = post.specs;
			post.n_links = post.links;
			
			post.detail = replaceAll('\r\n', '<br />', post.detail); //Replace those linebreaks
			post.specs = replaceAll('\r\n', '<br />', post.specs);
			post.links = replaceAll('\r\n', '<br />', post.links);
			
			var type = post.type;
			if (type == 'ee') {
				events.push(post);
			} else if (type == 'aa') {
				news.push(post);
			} else if (type == 'pp') {
				aircraft.push(post);
			} else if (type == 'vv') {
				vehicles.push(post);
			} else if (type == 'pr') {
				projects.push(post);
			}
		}
		
		//Sort events and news
		events.sort(function(a,b){
	  		return new Date(a.date) - new Date(b.date);
		});
		news.sort(function(a,b){
	  		return new Date(b.date) - new Date(a.date);
		});
		
		
		//console.log(events);
		var today = new Date();
		var li_offset = 0; // Offset for upcoming events list
		
		//For Home Page		
		$("#e-list").append('<li class="pure-menu-heading black bold left">Future Events</li>');
		for (var pst in events) {
			var post = events[pst];
			
			var title = post.title;
			var eday = new Date(post.date);
			if (today < eday) {
				$("#e-list").append('<li id="'+pst+'-0"><a href="javascript:setPost('+pst+',0)" class="pure-menu-link black left">'+title+'</a></li>');
			} else {
				li_offset = li_offset + 1;
			}
		}
		
		$("#e-list").append('<li class="pure-menu-heading black bold left">News</li>');
		for (var pst in news) {
			var post = news[pst];
			
			var title = post.title;
			
			$("#e-list").append('<li id="'+pst+'-1"><a href="javascript:setPost('+pst+',1)" class="pure-menu-link black left">'+title+'</a></li>');
		}
				
		setPost(li_offset,0);
		setCalendar();
		setArticles();
		setExhibits("pp");
		setUpEditor();
	});
	
});

function replaceAll(find, replace, str) {
  return str.replace(new RegExp(find, 'g'), replace);
}


var li_select = '#0-0'; //Selected event

function setPost(i,x) { //x = section HOMEPAGE
	
	var newID = '#'+i+'-'+x;
	$(li_select).removeClass("li-select");
	li_select = newID;
	$(li_select).addClass("li-select");
	
	//Update Info Box
	var post;
	if (x == 0) {
		post = events[i];
	} else if (x == 1) {
		post = news[i];
	}
	
	var title = document.getElementById("title");
	var date = document.getElementById("date");
	var details = document.getElementById("detail");
	
	if (title == null) return; //function not for this page
	
	var today = new Date(post.date);
	var dateString = today.format("mmmm d, yyyy");

	title.innerHTML = post.title;
	date.innerHTML = dateString;
	details.innerHTML = post.detail;
	
	
	var imgs = post.imgs;
	for (var i in imgs) {
		var img = imgs[i];
		$("#image").attr("src",img.path);
		
	}
	
}

function setCalendar() { // EVENTSNEWS
	
	var year = '1';
	for (var pst in events) {
		var post = events[pst];
		
		var date = new Date(post.date);
		var dateString = date.format("mmmm d");
		var yearString = date.format("yyyy");
		
		if (year != yearString) {
			$("#calendars").append('<h4 class="title">'+yearString+' Calendar</h4>');
			$("#calendars").append('<div id="c-list'+yearString+'" class="pure-g box rounded"></div>');
			year = yearString;
		}
		
		$('#c-list'+yearString).append('<a class="subTitle pure-u-1-2" href="/post/'+post.file+'">'+dateString+'</a> <p class="pure-u-1-2">'+post.title+'</p>');
		
	}
}

function setArticles() { // EVENTSNEWS

	for (var pst in news) {
		var post = news[pst];
		
		var date = new Date(post.date);
		var dateString = date.format("mmmm d, yyyy");
		
		//$('#p-list').append('<div id="'+post.file+'" class="box rounded"><a class="subTitle" href="/post/'+post.file+'">'+post.title+'</a> <p>'+dateString+'</p> <br /> <h4 class="subTitle">Details:</h4> <p>'+post.detail+'</p></div>');
		$('#p-list').append('<div id="" class="box rounded pure-g"><div class="pure-u-1 pure-u-md-1-2"><a class="subTitle" href="/post/'+post.file+'">'+post.title+'</a> <p>'+dateString+'</p> <br /> <h4 class="subTitle">Detail</h4> <p>'+post.detail+'</p> <br /> <h4 class="subTitle">Notes</h4> </div> <div id="'+post.file+'" class="pure-u-1 pure-u-md-1-2"> <h4 class="subTitle">Photos</h4> </div> </div>');
		appendPhotos(post,post.file)
	}
}


function setExhibits(exhibit) { // Exhibits
	$( "#exhibits" ).empty();
	
	$("#pp").removeClass("pure-menu-selected");
	$("#vv").removeClass("pure-menu-selected");
	$("#pr").removeClass("pure-menu-selected");
	
	$('#'+exhibit).addClass("pure-menu-selected");
	
	if (exhibit == 'pp') {
		for (var pst in aircraft) {
			var post = aircraft[pst];
		
			var date = new Date(post.date);
			var dateString = date.format("mmmm d, yyyy");
		
			$('#exhibits').append('<div id="" class="box rounded pure-g"><div class="pure-u-1 pure-u-md-1-2"><a class="subTitle" href="/post/'+post.file+'">'+post.title+'</a> <p>'+dateString+'</p> <br /> <h4 class="subTitle">Detail</h4> <p>'+post.detail+'</p> <br /> <h4 class="subTitle">Specs</h4> <p>'+post.specs+'</p> <br /> <h4 class="subTitle">Articles</h4> </div> <div id="'+post.file+'" class="pure-u-1 pure-u-md-1-2"> <h4 class="subTitle">Photos</h4> </div> </div>');
			appendPhotos(post,post.file)
		}
	} else if (exhibit == 'vv') {
		for (var pst in vehicles) {
			var post = vehicles[pst];
		
			var date = new Date(post.date);
			var dateString = date.format("mmmm d, yyyy");
		
			$('#exhibits').append('<div id="" class="box rounded pure-g"><div class="pure-u-1 pure-u-md-1-2"><a class="subTitle" href="/post/'+post.file+'">'+post.title+'</a> <p>'+dateString+'</p> <br /> <h4 class="subTitle">Detail</h4> <p>'+post.detail+'</p> <br /> <h4 class="subTitle">Specs</h4> <p>'+post.specs+'</p> <br /> <h4 class="subTitle">Articles</h4> </div> <div id="'+post.file+'" class="pure-u-1 pure-u-md-1-2"> <h4 class="subTitle">Photos</h4> </div> </div>');
			appendPhotos(post,post.file)
		}
	} else if (exhibit == 'pr') {
		for (var pst in projects) {
			var post = projects[pst];
		
			var date = new Date(post.date);
			var dateString = date.format("mmmm d, yyyy");
		
			$('#exhibits').append('<div id="" class="box rounded pure-g"><div class="pure-u-1 pure-u-md-1-2"><a class="subTitle" href="/post/'+post.file+'">'+post.title+'</a> <p>'+dateString+'</p> <br /> <h4 class="subTitle">Detail</h4> <p>'+post.detail+'</p> <br /> <h4 class="subTitle">Specs</h4> <p>'+post.specs+'</p> <br /> <h4 class="subTitle">Articles</h4> </div> <div id="'+post.file+'" class="pure-u-1 pure-u-md-1-2"> <h4 class="subTitle">Photos</h4> </div> </div>');
			appendPhotos(post,post.file)
		}
	}
}


function appendPhotos(pst,div) { //General Function

	var imgs = pst.imgs;
	var stackid = "imgstack"+div;
	$('#'+div).append('<div id="'+stackid+'" class="pure-g">');
	
	if (imgs.length > 2) { 
		for (var img in imgs) {
			var image = imgs[img];
			$('#'+stackid).append('<div class="pure-u-1 pure-u-sm-1-2 pure-u-md-1 pure-u-lg-1-2 l-img"> <img class="pure-img center" src="'+image.path+'"> </div>');
		}
	} else {
		for (var img in imgs) {
			var image = imgs[img];
			$('#'+stackid).append('<div class="pure-u-1 pure-u-sm-1-2 pure-u-md-1 pure-u-lg-1 l-img"> <img class="pure-img center" src="'+image.path+'"> </div>');
		}
	}
	
	$('#'+div).append('</div>');
}


function setUpEditor() { //Post Editor
	 
	$('#post-select').append('<optgroup label="Aircraft">');
	for (var pst in aircraft)  {
		var post = aircraft[pst];
		$('#post-select').append('<option id="0-'+pst+'">'+post.title+'</option>');
	}
	$('#post-select').append('</optgroup>');
	
	$('#post-select').append('<optgroup label="Vehicles">');
	for (var pst in vehicles)  {
		var post = vehicles[pst];
		$('#post-select').append('<option id="1-'+pst+'">'+post.title+'</option>');
	}
	$('#post-select').append('</optgroup>');
	
	$('#post-select').append('<optgroup label="Projects">');
	for (var pst in projects)  {
		var post = projects[pst];
		$('#post-select').append('<option id="2-'+pst+'">'+post.title+'</option>');
	}
	$('#post-select').append('</optgroup>');
	
	$('#post-select').append('<optgroup label="News">');
	for (var pst in news)  {
		var post = news[pst];
		$('#post-select').append('<option id="3-'+pst+'">'+post.title+'</option>');
	}
	$('#post-select').append('</optgroup>');
	
	$('#post-select').append('<optgroup label="Events">');
	for (var pst in events)  {
		var post = events[pst];
		$('#post-select').append('<option id="4-'+pst+'">'+post.title+'</option>');
	}
	$('#post-select').append('</optgroup>');
}

function getPostForEditor(id) {
	var res = id.split("-");
	var x = res[0];
	var y = res[1];
	
	if (x == 0) {
		return aircraft[y];
	} else if (x == 1) {
		return vehicles[y];
	} else if (x == 2) {
		return projects[y];
	} else if (x == 3) {
		return news[y];
	} else if (x == 4) {
		return events[y];
	}
}

