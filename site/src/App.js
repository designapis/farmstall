import React from 'react';
import { ReactComponent as FarmstallApiV1 } from './img/farmstallapi-v1.svg';
import BookCover from './img/Ponelat-Designing-MEAP-HI.jpg'
import { ReactComponent as LandingIllustration } from './img/landing-illustration.svg'
import './App.css';

const MANNING_LINK = "https://www.manning.com/books/designing-apis-with-swagger-and-openapi"

const SHUB_LINK = "https://app.swaggerhub.com/search?type=API&owner=designing-apis"

function App() {
  return (
    <div id="App">
      <header id="header">


        <FarmstallApiV1 id="logo" />
        <div id="blurb">
          <LandingIllustration id="landing-illustration" />

          <p>
            <h2>What is this?</h2>
            You're seeing the landing page of the FarmStall API.
            <br/>
            <br/>
            This RESTful API is an accompaniment to the book,
            <a href={MANNING_LINK}>
              Designing APIs with Swagger and OpenAPI
            </a>
            which is currently in Manning Early Access Program ( MEAP ).
          </p>

          <a href={MANNING_LINK}>
            <img
              id="book-cover"
              alt="Manning's Designing Apis with Swagger and OpenAPI. By Josh Ponelat"
              src={BookCover}/>
          </a>
        </div>

      </header>

      <div id="content">
        <div id="left-content">

          <h2> About the API </h2>

          <ul>
            <li> The API is hosted on https://farmstall.designapis.com/v1 </li>
            <li> The documentation of the API is a series of exercises  within the book, written as an OpenAPI (3.x) definition and described using Swagger tooling. </li>
            <li>
              The stages of the API definition can be found in <a href={SHUB_LINK}> SwaggerHub ( where I work ) </a>
            </li>
          </ul>
        </div>
      </div>

    </div>
  );
}

export default App;
