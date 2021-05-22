<script>
    import imageless_product_svg from '../../../resources/product.svg';
    import FieldData, { FieldStates } from '../../../classes/FieldData';
    import Button from '../../../components/Button.svelte';
    import Input from '../../../components/Input.svelte';
    import { server_name } from '../../../server_info';
    import { onMount } from 'svelte';

    export let close = () => {};
    export let session_key = "";
    export let product_data = {};

    let product_images = [];
    let current_image_index = 0;
    let purchase_amount = new FieldData("purchase-amount-input", /\d+/, "count", "number");

    onMount(() => { 
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const request = new Request(`http://${server_name}/products-images?id=${product_data.id}`, {method: 'GET', headers: headers});
        fetch(request)
            .then(promise => {
                if (promise.ok) {
                    promise.json().then(images => product_images = images);
                }
            });
    });

    const changeImage = direction => {
        switch (direction) {
            case "+":
                current_image_index = current_image_index == (product_images.length - 1) ? 0 : current_image_index + 1;
                break;
            case "-":
                current_image_index = current_image_index > 0 ?  current_image_index - 1 : product_images.length - 1;
            default:
                break;
        }
    }

    const handleBackgroundClicked = e => {
        if(e.currentTarget === e.target) {
            close();
        }
    }

    const restrinAmount = e => {
        const { target:count_element } = e;
        let count_value = parseInt(count_element.value);
        if (count_value > (product_data.stock - product_data.solds)) {
            count_value = (product_data.stock - product_data.solds)
        } else if (count_value <= 0) {
            count_value = 1;
        }
        count_element.value = count_value;
    }

    const makeOrder = () => {
        const headers = new Headers();
        headers.set("X-sk", session_key);

        const form_data = new FormData();
        form_data.append("product_id", product_data.id);
        form_data.append("count", purchase_amount.getFieldValue());

        const request = new Request(`http://${server_name}/cart`, {method:'POST', body: form_data, headers: headers});
        fetch(request).then(() => close());
    }

</script>

<style>
    #buy-product-modal-background {
        position: fixed;
        display: flex;
        background-color: var(--modal-shadows);
        width: 100vw;
        height: 100vh;
        justify-content: center;
        align-items: center;
    }

    #bp-modal {
        background: white;
        width: 60%;
        height: 70vh;
        padding: 1vh 1vw;
        border-radius: 5px;
    }

    #product-information {
        height: 60vh;
    }

    #product-information-upper {
        display: grid;
        height: 37vh;
        grid-template-columns: 1.2fr 2fr;
    }

    #product-image-container {
        display: flex;
        width: 100%;
        justify-content: center;
        align-items: center;
    }

    #product-image-container img {
        max-width: 13vw;
        height: calc(inherit - .2vh);
        border-radius: 8px;
    }

    .image-changer-btn {
        cursor: pointer;
        width: 4vw;
        height: 37vh;
    }

    #imageless-logo-container {
        width: 10vw;
    }

    #product-images {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    #product-general {
        display: flex;
        flex-direction: column;
    }

    #bp-product-name {
        width: 100%;
        background: var(--theme-color);
        color: white;
        font-size: 1.3rem;
        text-align: center;
        text-transform: uppercase;
        margin: 2vh 0 6vh 0;
    }

    #pricing-data {
        display: flex;
        font-size: 1.2rem;
        padding-left: 4vw;
        flex-direction: column;
    }

    #product-description-label {
        text-transform: capitalize;
        font-size: 1.4rem;
        border-bottom: 2px solid var(--theme-color);
        padding: 0 2vw;
    }

    #bp-product-description {
        
        padding: 1vh 1vw;
        font-size: 1.1rem;
    }

    #bp-controls {
        display: flex;
        padding: 0 2vw;
        justify-content: flex-end;
        align-items: center;
    }

    #purchase-amount-container {
        margin-right: 2vw;
    }
</style>

<div on:click={handleBackgroundClicked} id="buy-product-modal-background">
    <div id="bp-modal">
        <div id="product-information">
            <div id="product-information-upper">
                <div id="product-images">
                    {#if product_images.length > 0}
                        <div id="product-image-container">
                            <div class="image-changer-btn" on:click={() => changeImage("-")}></div>
                            <img src="http://{server_name}/static/{product_images[current_image_index]}" alt="">
                            <div class="image-changer-btn" on:click={() => changeImage("+")}></div>
                        </div>
                    {:else}
                        <div id="imageless-logo-container">
                            {@html imageless_product_svg}
                        </div>
                    {/if}
                </div>
                <div id="product-general">
                    <span id="bp-product-name">{product_data.name}</span>
                    <div id="pricing-data">
                        <span id="bp-product-price">Price: {product_data.price.toLocaleString('en')} MXN</span>
                        <span id="bp-product-">Stock: {product_data.stock - product_data.solds}</span>
                    </div>
                </div>
            </div>
            <div id="product-information-lower">
                <div id="product-description">
                    <div id="product-description-label">
                        Descripcion
                    </div>
                    <span id="bp-product-description">{product_data.description}</span>
                </div>
            </div>
        </div>
        <div id="bp-controls">
            <div id="purchase-amount-container">
                <Input initial_value="1" onBlur={restrinAmount} input_label="count" input_padding=".1vh .1vw" field_data={purchase_amount} min={0} max={product_data.stock - product_data.solds}/>
            </div>
            <Button onClick={makeOrder} width="4vw" label="purchase" />
        </div>
    </div>
</div>