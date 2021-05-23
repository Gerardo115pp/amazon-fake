<script>
    import { FieldStates } from '../../classes/FieldData'
    import Button from '../../components/Button.svelte';
    import Input from '../../components/Input.svelte';
    import FieldData from '../../classes/FieldData';
    import { server_name } from '../../server_info';
    import { push } from 'svelte-spa-router';
    
    let form_inputs = [
        new FieldData("newuser-username", /^[a-z_\d]+$/g, "username"),
        new FieldData("newuser-name", /^[A-Za-z_\s\d]+$/g, "name"),
        new FieldData("newuser-phone", /^[\d-]+$/g, 'phone',"tel"),
        new FieldData("newuser-email", /^[a-zA-Z\d_\.]+@[a-zA-Z]+(.[a-zA-Z])+$/g, "email","email"),
        new FieldData("newuser-address", /^[a-z_\s#\d]+$/g, "address"),
        new FieldData("newuser-password", /[^\s\n]+/g, "password","password")
    ];

    const awaitEnter = e => {
        if(e.key.toLowerCase() === "enter") {
            registerUser();
        }
    }

    const checkForm = () => {
        let is_form_ready = true;

        for(let fd of form_inputs) {
            if(!fd.isReady()) {
                fd.state = FieldStates.HAS_ERRORS;
                is_form_ready = false;
            } else {
                fd.state = FieldStates.READY;
            }
        }
        form_inputs = [...form_inputs];
        return is_form_ready;
    }

    const registerUser = () => {
        if(checkForm()) {
            const form_data = new FormData();
            for(let fd of form_inputs) {
                form_data.append(fd.name, fd.getFieldValue());
            }

            const request = new Request(`http://${server_name}/register`, { method: 'POST', body: form_data});
            fetch(request)
                .then(promise => {
                    if (promise.ok) {
                        push("/");
                    } else {
                        promise.json().then(response => alert(response.error))
                    }
                })
        }
    }

</script>

<style>
    #app-sign-in-page {
        height: 100vh;
    }

    #register-title {
        margin: 5vh 0;
        text-align: center;
    }

    #register-form-container {
        display: flex;
        width: 80%;
        flex-direction: column;
        margin: 10vh auto 5vh;
    }

    .register-field-container {
        box-sizing: border-box;
        height: 5vh;
        margin-bottom: 5vh;
    }

    .controls {
        display: flex;
        width: 40%;
        justify-content: center;
        align-self: center;

    }
</style>

<div id="app-sign-in-page">
    <h1 id="register-title">Registro</h1>
    <div id="register-form-container">
        {#each form_inputs as fd}
            <div class="register-field-container">
                <Input
                    field_data={fd}
                    input_padding=".5vh 2.5vw"
                />
            </div>
        {/each}
        <div class="controls">
            <Button  onClick={registerUser} width="20%" label="Registrar"/>
        </div>
    </div>
</div>