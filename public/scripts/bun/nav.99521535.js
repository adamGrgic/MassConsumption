// @bun
console.log("nav js loaded");var l=document.querySelectorAll("#navbar a");console.log(l);l.forEach((e)=>{console.log("adding event listener for link: ",e),e.addEventListener("click",(o)=>{console.log("click"),o.preventDefault(),console.log(`Clicked link: ${e}`),console.log(o),e.setAttribute("active","true")})});
