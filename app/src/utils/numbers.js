export const parseValue = (value) => {
    const val = parseFloat(value);
    return isNaN(val) ? null : val;
};