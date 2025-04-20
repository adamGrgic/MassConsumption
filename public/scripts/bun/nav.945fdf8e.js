// @bun
console.log("nav js loaded");var o=document.querySelectorAll("#navbar a");console.log(o);o.forEach((e)=>{console.log("adding event listener for link: ",e),e.addEventListener("click",(l)=>{console.log("click"),l.preventDefault(),console.log(`Clicked link: ${e}`),console.log(l),o.forEach((c)=>c.removeAttribute("active")),e.setAttribute("active","true")})});
