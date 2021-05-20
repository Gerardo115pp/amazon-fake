export const FieldStates = {
    NORMAL: 0,
    HAS_ERRORS: 1,
    READY: 2
}

class FieldData {
    constructor(field_id, validation_regex, name,type_name="text") {
        this.id = field_id;
        this.name = name;   
        this.regex = validation_regex;
        this.type = type_name;
        this.state = FieldStates.NORMAL;

    }

    clear = () => {
        this.getField().value = '';
    }

    getField = () => {
        return document.getElementById(this.id);
    }

    getFieldValue = () => {
        return this.getField().value;
    }

    isReady = () => {
        return this.regex.test(this.getFieldValue());
    }
}

export default FieldData;