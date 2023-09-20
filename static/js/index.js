const arrayContainerProdutos = document.querySelectorAll('.outer-container')
const arrayProdutos = document.querySelectorAll('.inner-container')
const arrowImage = document.querySelectorAll('.arrow-img')

arrayContainerProdutos.forEach((element, index) => {
   element.addEventListener('click', ()=>{
    arrayProdutos[index].classList.toggle('hidden')
    arrowImage[index].classList.toggle('before:rotate-[135deg]')
    arrowImage[index].classList.toggle('before:-rotate-45')
    arrowImage[index].classList.toggle('before:top-[30px]') 
    arrowImage[index].classList.toggle('before:top-[45px]') 
   }) 
});