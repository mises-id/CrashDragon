var currentProduct = document.cookie.replace(/(?:(?:^|.*;\s*)product\s*\=\s*([^;]*).*$)|^.*$/, "$1");
if (currentProduct != "all") {
    $("#product-filter").val(currentProduct);
}

$("#product-filter").on("change", function () {
    var slug = $("#product-filter").val();
    document.cookie = "product=" + slug + ";path=/";
    location.reload();
});

$('#change-slug').on("click", function() {
    alert("This can have some unintended side-effects! Be sure to know what you do!");
    var slugField = document.getElementById('slug');
    if (slugField.hasAttribute('readonly')) {
        slugField.removeAttribute('readonly');
        $('#change-slug').hide();
    }
});
