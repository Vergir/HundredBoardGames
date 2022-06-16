import * as Lang from "/static/js/lang.js";

function handleGameMoreButton(event) {
    const buttonElement = this;
    const gameElement = buttonElement.closest('.game');

    toggleGame(gameElement);
}

function toggleGame(gameElement) {
    const isExpanded = gameElement.classList.contains("game--big");
    if (isExpanded) {
        minimizeGameCard(gameElement);
    } else {
        maximizeGameCard(gameElement);
    }
}

function IsSkipKeyboardEvent(event, targetElement) {
    if (!(event instanceof KeyboardEvent)) {
        return false
    }

    if (targetElement && event.target !== targetElement) {
        return true;
    }

    const validKeys = ['Space', 'Enter'];
    
    const needSkip = !validKeys.includes(event.code); 
    
    return needSkip;
}

function maximizeGameCard(gameElement) {
    gameElement.classList.add("game--big");
    
    gameElement.querySelector(".gameExtraBorder").classList.remove('hidden');
    gameElement.querySelector(".gameExtra").classList.remove('hidden');
    gameElement.querySelector('.gameExpandArrow').classList.add('gameExpandArrow--reverse');
    gameElement.querySelectorAll(".gameExtraPictures lazy").forEach(pictureElement => pictureElement.classList.add('lazypreload'));
    
    gameElement.scrollIntoView();
}

function minimizeGameCard(gameElement) {
    gameElement.classList.remove('game--big');
    
    gameElement.querySelector('.gameExtraBorder').classList.add('hidden');
    gameElement.querySelector('.gameExtra').classList.add('hidden');
    gameElement.querySelector('.gameExpandArrow').classList.remove('gameExpandArrow--reverse');

    const needScroll = gameElement.getBoundingClientRect().top < 0;
    if (needScroll) {
        gameElement.scrollIntoView();
    }
}

function handleSmallPictureKbm(event) {
    const smallPictureElement = this;
    
    const picId = smallPictureElement.dataset.picId;

    const smallPicturesList = smallPictureElement.closest('.gameExtraSmallPictures');

    activateSmallPicture(smallPicturesList, picId, smallPictureElement);

    const bigPicturesList = smallPicturesList.closest('.gameExtra').querySelector('.gameExtraBigPictures');
    const targetBigPicture = bigPicturesList.querySelector(`[data-pic-id="${picId}"]`);
    bigPicturesList.scrollLeft = targetBigPicture.offsetLeft - (bigPicturesList.offsetWidth/2);
}

function activateSmallPicture(smallPicturesList, picId, smallPictureElement) {    
    smallPictureElement = smallPictureElement || smallPicturesList.querySelector(`[data-pic-id="${picId}"]`);

    const activeClass = 'gameExtraSmallPicture--active';
    smallPicturesList.querySelector('.' + activeClass).classList.remove(activeClass);
    smallPictureElement.classList.add(activeClass);

    const rem = parseFloat(getComputedStyle(document.documentElement).fontSize);

    smallPicturesList.scrollLeft = smallPictureElement.offsetLeft - smallPicturesList.offsetWidth/2 + smallPictureElement.offsetWidth/2;
}

let scrollTimerId = -1;
function handleBigPicturesListScroll(event) {
    const bigPicturesList = event.target;

    if (scrollTimerId != -1) {
        clearTimeout(scrollTimerId);
    }

    scrollTimerId = window.setTimeout(scrollFinished, 10, bigPicturesList);      
}

function scrollFinished(bigPicturesList) {
    let targetBigPicture = null;
    bigPicturesList.childNodes.forEach((bigPictureElement) => {
        if (targetBigPicture) {
            return;
        }
        if (bigPictureElement.offsetLeft >= (bigPicturesList.scrollLeft - bigPicturesList.offsetWidth/2)) {
            targetBigPicture = bigPictureElement;
        }
    });

    const smallPicturesList = bigPicturesList.closest('.gameExtra').querySelector('.gameExtraSmallPictures');

    activateSmallPicture(smallPicturesList, targetBigPicture.dataset.picId);
}

function handleGameExtraHideButtonClick(event) {
    const buttonElement = this;
    const gameElement = this.closest('.game');

    const needToShiftFocus = (document.activeElement === buttonElement)
    if (needToShiftFocus) {
        this.closest('.game').querySelector('.gameMoreButton').focus();
    }

    minimizeGameCard(gameElement);
}

function showAttributeTooltip(gameAttributeElement) {
    const gameElement = gameAttributeElement.closest('.game');

    let tooltipPositionClass = "tooltip--above";
    if (gameElement.getBoundingClientRect().top <= 200) {
        tooltipPositionClass = "tooltip--below";
    }

    const content = getTooltipContent(gameAttributeElement);

    const html = `<div class="tooltip ${tooltipPositionClass}">${content}</div>`;

    gameAttributeElement.insertAdjacentHTML('afterend', html);

    const stylesheet = document.styleSheets[0];
    const hasCustomOffsetRule = stylesheet.cssRules[stylesheet.cssRules.length-1].selectorText.startsWith('.tooltip');
    if (hasCustomOffsetRule) {
        stylesheet.deleteRule(stylesheet.cssRules.length-1);
    }
    const offset = gameAttributeElement.offsetLeft + gameAttributeElement.offsetWidth/2;
    stylesheet.insertRule(`.${tooltipPositionClass}::after { left: calc(${offset}px - var(--arrow-size) - var(--scaling-factor)); }`, stylesheet.cssRules.length);
}

function getTooltipContent(element) {
    let content = '';
    const gameElement = element.closest('.game');
    const elementTypeClass = element.classList[element.classList.length-1];
    const isAttributeWithTooltip = ['gameComplexity', 'gameNumPlayers', 'gameLength'].includes(elementTypeClass);
    if (!isAttributeWithTooltip) {
        return '...';
    }

    const badge = `<span class="gameMainAttribute ${elementTypeClass}">` + Lang.get(`badge_${elementTypeClass}`) + '</span>';
    const commonDescription = Lang.expand(Lang.get(`description_common_${elementTypeClass}`), {'badge': badge});

    let specificDescriptionTokens = {
        'game': `<i>${gameElement.dataset.title}</i>`,
        'attribute': element.outerHTML.replace('tabindex="0"', ''),
    };
    switch (elementTypeClass) {
        case 'gameComplexity':
            specificDescriptionTokens = {
                ...specificDescriptionTokens,
                'weight': '<b>' + (parseFloat(gameElement.dataset.avgWeight) * 2).toFixed(1) + '&nbsp;/&nbsp;10.0</b>',
                'weight_num': `<b>${gameElement.dataset.weightNumVotes}</b>`,
            };
            break;
        case 'gameNumPlayers':
            const gameNumSpecificDescriptionTokens = formNumPlayersSpecificTooltipTokens(gameElement);
            specificDescriptionTokens = {
                ...specificDescriptionTokens,
                ...gameNumSpecificDescriptionTokens,
            };
            break;
    }
    const specificDescription = Lang.expand(Lang.get(`description_specific_${elementTypeClass}`), specificDescriptionTokens);
    const label = `<p>${commonDescription}</p><p>${specificDescription}</p>`;

    return label;
}

function formNumPlayersSpecificTooltipTokens(gameElement) {
    const numPlayers = gameElement.querySelector('.gameNumPlayers').innerText;
    let [minPlayers, maxPlayers] = [0, 0];
    if (numPlayers.includes('-')) {
        [minPlayers, maxPlayers] = numPlayers.split('-').map(Number);
    } else {
        [minPlayers, maxPlayers] = [Number(numPlayers), Number(numPlayers)];
    }

    const communityNumPlayersDataset = gameElement.querySelector('.communityNumPlayers').dataset;
    let bestNumPlayers = 0;
    let bestNumPlayersVotes = 0;
    let communityMinPlayers = 0;
    let communityMaxPlayers = 0;
    Object.values(communityNumPlayersDataset).forEach((codedData) => {
        const [numPlayers, votedBest, votedRecommended, votedNotRecommended] = codedData.split(',').map(Number);
        if (votedBest > bestNumPlayersVotes) {
            bestNumPlayers = numPlayers;
            bestNumPlayersVotes = votedBest;
        }

        const hasCommunityMinPlayers = (numPlayers < minPlayers) && (votedRecommended > votedNotRecommended);
        const hasCommunityMaxPlayers = (numPlayers > maxPlayers) && (votedRecommended > votedNotRecommended);
        if (hasCommunityMinPlayers) {
            communityMinPlayers = numPlayers
        }
        if (hasCommunityMaxPlayers) {
            communityMaxPlayers = numPlayers
        }
    });

    const communityPlayersClauses = [];
    if (communityMinPlayers) {
        const peopleLabel = '<b>' + Lang.numPeople(communityMinPlayers) + '</b>';
        const minPlayersClause = Lang.expand(Lang.get('community_num_players_clause_min'), {'num_people': peopleLabel});
        communityPlayersClauses.push(minPlayersClause);
    }
    if (communityMaxPlayers) {
        const peopleLabel = '<b>' + Lang.numPeople(communityMaxPlayers) + '</b>';
        const maxPlayersClause = Lang.expand(Lang.get('community_num_players_clause_max'), {'num_people': peopleLabel});
        communityPlayersClauses.push(maxPlayersClause);
    }
    let communityNumPlayers = '';
    if (communityPlayersClauses.length) {
        const communityNumPlayersTokens = {
            'clauses': communityPlayersClauses.join(Lang.get('community_num_players_clauses_link')),
        };
        communityNumPlayers = Lang.expand(Lang.get('community_num_players'), communityNumPlayersTokens);
    }

    const tokens = {
        'community_num_players': communityNumPlayers,
        'num_people': '<b>' + Lang.numPeople(bestNumPlayers) + '</b>',
    }

    return tokens;
}


function hideTooltips() {
    [...document.getElementsByClassName('tooltip')].forEach(
        tooltipElement => tooltipElement.remove()
    );
}

function handleGameMainAttributeKbm(event) {
    const gameMainAttributeElement = this;

    if (IsSkipKeyboardEvent(event, gameMainAttributeElement)) {
        return;
    }

    event.preventDefault();

    hideTooltips();
    
    showAttributeTooltip(gameMainAttributeElement);
}

function initGame(gameElement) {
    gameElement.dataset.title = gameElement.querySelector('.gameTitle').innerText;
}

function handleGameKeydown(event) {
    const gameElement = this;

    if (IsSkipKeyboardEvent(event, gameElement)) {
        return;
    }

    event.preventDefault(); //so pressing Space doesn't scroll.

    toggleGame(gameElement);
}

function handleBodyClick(event) {
    const clickInsideTooltip = (event.target.closest('.tooltip') != null);
    const clickToShowTooltip = (event.target.closest('.gameMainAttribute') != null);
    const needToHideTooltips = !clickInsideTooltip && !clickToShowTooltip;
    if (needToHideTooltips) {
        hideTooltips();
    }
}

function handleBodyKeydown(event) {
    if (event.code != "Tab") {
        return;
    }

    hideTooltips();
}

function addHandlers(selector, handler, ...events) {
    document.querySelectorAll(`${selector}`).forEach((element) => {
        events.forEach(eventName => element.addEventListener(eventName, handler));
    });
}

/* === */

document.querySelectorAll('.game').forEach(
    gameElement => initGame(gameElement)
);

addHandlers('body', handleBodyClick, 'click');
addHandlers('body', handleBodyKeydown, 'keydown');
addHandlers('.game', handleGameKeydown, 'keydown');
addHandlers('.gameMoreButton', handleGameMoreButton, 'click');
addHandlers('.gameExtraSmallPicture', handleSmallPictureKbm, 'click', 'keydown');
addHandlers('.gameExtraBigPictures', handleBigPicturesListScroll, 'scroll');
addHandlers('.gameExtraHideButton', handleGameExtraHideButtonClick, 'click');
addHandlers('.gameMainAttribute', handleGameMainAttributeKbm, 'click', 'keydown');
