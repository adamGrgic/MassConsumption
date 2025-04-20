// @bun
console.log("nav js loaded");var e=document.querySelectorAll("#navbar a");console.log(e);e.forEach((o)=>{console.log("adding event listener for link: ",o),o.addEventListener("click",(l)=>{console.log("click"),l.preventDefault(),console.log(`Clicked link: ${o}`),console.log(l)})});
