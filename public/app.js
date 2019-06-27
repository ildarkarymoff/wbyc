const State = {
    IDLE: 0,
    FETCHING: 1,
    FETCH_COMPLETE: 2,
    ERROR: 3
};

const WeatherApp = {
    init: () => {
        WeatherApp.state = State.IDLE;

        WeatherApp.weather = null;
        WeatherApp.API_URI = '/api/weather/current/';
        WeatherApp.UI = {
            inputField: document.querySelector('.cityInput'),
            showButton: document.querySelector('.showBtn'),
            output: document.querySelector('.output')
        };

        WeatherApp.UI.showButton.addEventListener('click', WeatherApp.fetchWeather);
        WeatherApp.UI.inputField.addEventListener('keyup', e => {
            if (e.key === 'Enter') WeatherApp.fetchWeather();
        });
    },
    updateView: error => {
        switch (WeatherApp.state) {
            case State.FETCHING:
                WeatherApp.UI.output.innerHTML = '<i>Please wait...</i>';
                break;
            case State.FETCH_COMPLETE:
                WeatherApp.UI.output.innerHTML = `<p>Now in <strong>${WeatherApp.city}:</strong></p>`
                    + `<p class='large'>${WeatherApp.weather.temperature} &#8451;</p>`
                    + `<p class='feelsLike'>Feels like <strong>${WeatherApp.weather.feelsLike} &#8451;</strong></p>`;
                break;
            case State.ERROR:
                WeatherApp.UI.output.innerHTML = `<p class="error">${error}</p>`;
                break;
        }
    },
    fetchWeather: () => {
        WeatherApp.state = State.FETCHING;
        WeatherApp.updateView();

        WeatherApp.city = WeatherApp.UI.inputField.value.trim();
        const url = `${WeatherApp.API_URI}/${WeatherApp.city}`;

        try {
            fetch(url)
                .then(res => {
                    if (res.status === 404) {
                        console.log('Not found!');
                        WeatherApp.state = State.ERROR;
                        WeatherApp.updateView('City not found.');
                        throw new Error('City not found');
                    }

                    return res.json();
                })
                .then(json => {
                    console.log(json);
                    if (json) {
                        if (json.temperature && json.feels_like)
                            WeatherApp.weather = {
                                temperature: json.temperature,
                                feelsLike: json.feels_like
                            };

                        WeatherApp.state = State.FETCH_COMPLETE;
                        WeatherApp.updateView();
                    } else throw new Error('Something wrong.');
                })
                .catch(e => {
                    console.log(e);
                    WeatherApp.state = State.ERROR;
                    WeatherApp.updateView(e);
                });
        } catch (e) {
            console.error(e);
            WeatherApp.state = State.ERROR;
            WeatherApp.updateView(e);
        }
    }
};

WeatherApp.init();