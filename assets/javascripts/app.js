var currentProduct = document.cookie.replace(/(?:(?:^|.*;\s*)product\s*\=\s*([^;]*).*$)|^.*$/, "$1");
if (currentProduct != "all") {
    $("#product-filter").val(currentProduct);
}

$("#product-filter").on("change", function () {
    var slug = $("#product-filter").val();
    document.cookie = "product=" + slug + ";path=/";
    location.reload();
});
