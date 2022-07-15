function get(key) {
    return lang[key] || `{${key}}`;
}

function expand(label, tokensExpansions) {
    Object.keys(tokensExpansions).forEach((token) => {
        label = label.replaceAll(`{${token}}`, tokensExpansions[token]);
    });

    return label;
}

function numPeople(num) {
    const peopleLabel = num === 1 ? get('person_1') : expand(get('person_num'), {'num': num});

    return peopleLabel;
}

const lang = {
    'show_more_info': 'Show more info',
    'hide_more_info': 'Hide more info',
    'badge_gameComplexity': 'Complexity',
    'description_common_gameComplexity': '{badge} measures how complex or difficult to understand the game is according to people who played it.',
    'description_specific_gameComplexity': '{game} has an average complexity {weight}, rated by {weight_num}&nbsp;people, thus marking it as {attribute}',
    'badge_gameLength': 'Game length',
    'description_common_gameLength': '{badge} shows how much time a typical playthrough takes.',
    'description_specific_gameLength': '{game} usually takes around {attribute}',
    'badge_gameNumPlayers': 'Number of players',
    'description_common_gameNumPlayers': '{badge} shows how many people can play the game at a time.',
    'description_specific_gameNumPlayers': '{game} oficially supports {attribute} people.\n{community_num_players}Most players agree that this game is best played by {num_people}.',
    'community_num_players': 'Conversely, the majority of players report the game can be played by as {clauses}. ',
    'community_num_players_clause_min': 'few as {num_people}',
    'community_num_players_clause_max': 'many as {num_people}',
    'community_num_players_clauses_link': ' and as ',

    'person_1': '1&nbsp;person',
    'person_num': '{num}&nbsp;people',
};

export {get, expand, numPeople};