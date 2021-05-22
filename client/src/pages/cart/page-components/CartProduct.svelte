<script>
    import remove_svg from '../../../resources/remove-from-cart.svg';
    import { server_name } from "../../../server_info";

    export let updateCart = () => {};
    export let session_key = "";
    export let product_data = {};

    const removeProduct = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const form_data = new FormData();
        form_data.append("id", product_data.product_id)

        const request = new Request(`http://${server_name}/cart`, {method: 'DELETE', body: form_data, headers: headers});
        fetch(request)
            .then(promise => {
                if(promise.ok) {
                    promise.json().then(new_cart => updateCart(new_cart));
                }
            })
    }
</script>

<style>
    .cart-product-container {
        cursor: pointer;
        display: flex;
        height: 10vh;
        background: rgba(0, 0, 0, 0.3);
        border-radius: 5px;
        margin-top: 1vh;
        justify-content: space-around;
        align-items: center;
        transition: all .2s ease-in;
    }

    .cart-product-container:hover {
        background: var(--ready);
    }

    .product-thumbnail-container {
        display: flex;
        align-items: center;
    }

    .product-thumbnail-container img{
        width: 4vw;
        height: 8vh;
        border-radius: 15px;
    }

    .product-data {
        color: #fff;
        display: flex;
        width: 30vw;
        align-items: center;
    }

    .product-data span:first-child {
        width: 70%;
    }

    .cart-product-remove-btn {
        width: 2vw;
        display: flex;
        align-items: center;
        transition: all .3s ease-in-out;
    }

    .cart-product-remove-btn:hover {
        transform: scale(1.2);
    }

</style>

<div class="cart-product-container">
    <div class="product-thumbnail-container">
        <img src="http://{server_name}/static/{product_data.thumbnail}" alt="{product_data.thumbnail}">
    </div>
    <div class="product-data">
        <span class="cart-product-name">{product_data.product_name}</span>
        <span class="cart-product-total">${product_data.unit_price.toLocaleString('en')} MXN</span>
    </div>
    <div on:click={removeProduct} class="cart-product-remove-btn">
        {@html remove_svg}
    </div>
</div>