<script>
    import FieldData from '../classes/FieldData';
    import { FieldStates } from '../classes/FieldData';

    export let field_data = new FieldData("generic-input", /[.\n]+/, "any"); 
    export let input_label;
    export let input_padding = "1.5vh 4.2vh";
    export let onKeypressed;
    export let onBlur;

    export let initial_value;
    export let min = 0;
    export let max;

    $: state_color = getBorderColor(field_data.state) ;


    const handleOutsideClickDetected = e => {
        if(e.currentTarget === e.target) {
            e.target.getElementsByTagName("input")[0].focus()
        }
    }

    const getBorderColor = state => {
        switch(state) {
            case FieldStates.NORMAL:
                return "--theme-color";
            case FieldStates.HAS_ERRORS:
                return "--danger";
            case FieldStates.READY:
                return "--ready";
            default:
                return "--theme-color"
        }
    }
</script>

<style>
    .input-container {
        cursor: text;
        display: flex;
        border: 2px solid var(--theme-color);
        border-radius: 9999px;
        align-items: center;
    }   

    .input-container label {
        margin-left: 1vw;
    }

    .input-container input {
        font-family: var(--main-font);
        width: 100%;
        display: flex;
        background: none;
        border: none;
        color: var(--theme-color);
        font-size: 1.2rem;
        padding: 0;
        align-items: center;
        outline: none;
    }

    .input-container input[type='number'] {
        width: 5rem;
        text-align: center;
        -moz-appearance: textfield;
    }

    .input-container input::placeholder {
        font-family: 'Rhodium Libre';
        color: var(--dimtheme-color);
        text-transform: capitalize;
    }

    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }


</style>

<div on:click={handleOutsideClickDetected} style="padding: {input_padding};border-color: var({state_color});" class="input-container">
    {#if input_label !== undefined}
        <label for={field_data.id}>{input_label}:</label>
    {/if}
    {#if field_data.type !== "number"}
        <input id={field_data.id} type={field_data.type}
            placeholder={field_data.name}
            on:keydown={onKeypressed}
            on:blur={onBlur}
        />
    {:else}
        <input id={field_data.id} type={field_data.type}
            placeholder={field_data.name}
            on:keydown={onKeypressed}
            on:blur={onBlur}
            value={initial_value}
            min={min}
            max={max}
        />
    {/if}

</div>