// @bun
console.log("nav js loaded");var l=document.querySelectorAll("#navbar a");console.log(l);l.forEach((o)=>{console.log("adding event listener for link: ",o),o.addEventListener("click",(e)=>{console.log("click"),e.preventDefault(),console.log(`Clicked link: ${o}`),console.log(`event ${e}`)})});
