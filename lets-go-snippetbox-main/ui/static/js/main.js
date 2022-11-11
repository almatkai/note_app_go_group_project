var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}
var active = document.querySelectorAll(".note");
active[0].classList.add("live1")
active[0].firstElementChild.classList.add("live2");
// active[0].firstElementChild.classList.add("ab")
// document.querySelectorAll("#list-notes").
// active[i].addEventListener("mouseover", function () {
// 	this.firstChild.classList.add("ab");
// });