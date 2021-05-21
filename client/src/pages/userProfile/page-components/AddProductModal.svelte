<script>
    import images_btn_svg from '../../../resources/images-btn.svg';
    import FieldData from '../../../classes/FieldData';
    import { FieldStates } from '../../../classes/FieldData';
    import Button from '../../../components/Button.svelte';
    import Input from '../../../components/Input.svelte';
    import { server_name } from '../../../server_info';


    export let onProductAdded = () => {};
    export let close = () => {};
    export let session_key = "";

    const files_input_id = "files-input";
    let product_images = [];

    let form_inputs = [
        new FieldData("newproduct-name", /^.+$/, "name"),
        new FieldData("newproduct-description", /^.+$/, "description"),
        new FieldData("newproduct-price", /^\d+(\.[\d]{1,3})?$/, "price"),
        new FieldData("newproduct-stock", /^\d+$/, "stock")
    ];


    const awaitEnter = e => {
        if(e.key.toLowerCase() === "enter") {
            // registerUser();
        }
    }

    const checkForm = () => {
        let is_form_ready = true;

        for(let fd of form_inputs) {
        if(fd.getFieldValue() === ""){
            fd.state = FieldStates.NORMAL;
            is_form_ready = false;  
        } else if(!fd.isReady()) {
                fd.state = FieldStates.HAS_ERRORS;
                is_form_ready = false;
        } else {
                fd.state = FieldStates.READY;
            }
        }
        form_inputs = [...form_inputs];
        return is_form_ready;
    }

    const handleBackgroundClick = e => {
        if (e.target === e.currentTarget) {
            close();
        }
    }

    const readFile = file => {
        return new Promise((onSuccess, onReject) => {
            let file_reader = new FileReader();
            file_reader.onload = () => onSuccess({src: file_reader.result, name: file.name, blob: file});
            file_reader.onerror = () => onReject(fr);
            file_reader.readAsDataURL(file);
        });
    }

    const sumitProduct = () => {
        if(checkForm()) {
            const headers = new Headers();
            headers.set("X-sk", session_key);
            
            const form_data = new FormData();
            form_inputs.forEach(fd => form_data.append(fd.name, fd.getFieldValue()));
            if (product_images.length > 0) {
                product_images.forEach((images_data, h) => {
                    form_data.append(`image-${h}`, images_data.blob, images_data.name);
                });
            }

            const request = new Request(`http://${server_name}/products`, { method: 'POST', body: form_data, headers: headers});
            fetch(request)
                .then(promise => {
                    if(promise.ok) {
                        onProductAdded();
                        close();
                    }
                });
        }
    }

    const saveFiles = e => {
        const { target:file_input } = e;
        const promises = []

        for(let f of file_input.files) {
            if(/image.*/.test(f.type)) {
                promises.push(readFile(f))
            }
        }
        Promise.all(promises).then(files => product_images = files);
    }
    
    const triggerFilesInput = () => document.getElementById(files_input_id).click();
</script>

<style>
    #add-product-background {
        display: flex;
        position: fixed;
        width: 100vw;
        height: 100vh;
        background-color: rgba(0, 0, 0, 0.2);
        justify-content: center;
    }

    #ap-modal {
        width: 80vw;
        background-color: white;
        border-left: 3px solid var(--theme-color);
        border-right: 3px solid var(--theme-color);
        box-shadow: inset 0 0 60px 15px rgba(4, 205, 255, 0.137);
    }

    #ap-add-image-product {
        display: flex;
        background-color: white;
        height: 30vh;
        justify-content: space-around;
        align-items: center;
        border-bottom: 1px solid var(--theme-color);
    }

    #images-container {
        display: flex;
        justify-content: left;
        overflow-x: auto;
        width: 50vw;
        height: 20vh;
    }

    .product-image-container {
        display: flex;
        margin-right: 1vw;
        align-items: center;
    }

    .product-image-container img{
        width: 15vw;
        border-radius: 5px;
    }

    #add-product-images {
        cursor: pointer;
        display: flex;
        width: 6vw;
        background: var(--theme-gradiant);
        padding: .7vh .7vw;
        border-radius: 5px;
        justify-content: center;
        align-items: center;
    }

    :global(#add-product-images svg) {
        width: 4vw;
        fill: white;
        margin: 0  0 .5vw .5vw;
    }

    #new-product-form {
        display: flex;
        width: 80%;
        flex-direction: column;
        margin: 5% auto;
        align-items: center;
    }

    .new-product-field {
        width: 80%;
        margin-bottom: 4vh;
    }

    #files-mountpoint {
        display: none;
    }

</style>

<div on:click={handleBackgroundClick} id="add-product-background">
    <div id="ap-modal">
        <div id="ap-add-image-product">
            <div id="images-container">
                {#each product_images as image_data}
                    <div class="product-image-container">
                        <img src={image_data.src} alt="{image_data.name}">
                    </div>
                {/each}
            </div>
            <div on:click={triggerFilesInput} id="add-product-images">
                {@html images_btn_svg}
            </div>
        </div>
        <div id="new-product-form">
            {#each form_inputs as fd}
                <div class="new-product-field">
                    <Input onBlur={checkForm} onKeypressed={awaitEnter} input_padding=".3vh 2vw" field_data={fd} />
                </div>
            {/each}
            <div class="form-controls">
                <Button onClick={sumitProduct} width="15vw" label="Sumit product"/>
            </div>
            <div id="files-mountpoint">
                <input id={files_input_id} on:change={saveFiles} type="file" multiple/>
            </div>
        </div>
    </div>
</div>