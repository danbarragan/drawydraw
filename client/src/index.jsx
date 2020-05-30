import React from 'react';
import ReactDOM from 'react-dom';
import { IntlProvider } from 'react-intl';

import EnglishMessages from './translations/en.json';
import SpanishMessages from './translations/es.json';
import Game from './components/Game/Game';

const messages = {
  es: SpanishMessages,
  en: EnglishMessages,
};
// Split the browser's locale string to get the language without the region
const language = navigator.language.split(/[-_]/)[0];

ReactDOM.render(
  <IntlProvider locale={language} messages={messages[language]}><Game /></IntlProvider>,
  document.getElementById('root'),
);
