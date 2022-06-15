function get(key) {
    return lang[key] || `{${key}}`;
}

function expand(label, tokensExpansions) {
    Object.keys(tokensExpansions).forEach((token) => {
        label = label.replaceAll(`{${token}}`, tokensExpansions[token]);
    });

    return label;
}

const lang = {
    'badge_game_complexity': 'Complexity',
    'game_complexity_description_common': '{badge} measures how complex or difficult to understand the game is according to people who played it.',
    'game_complexity_description_specific': '{game} has an average complexity {weight}, rated by {weight_num} people, thus marking it as {complexity}'
};

export {get, expand};