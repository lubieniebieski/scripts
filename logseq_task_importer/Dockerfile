FROM ruby:3-alpine

RUN apk update && apk add build-base

RUN mkdir /script

WORKDIR /script

COPY Gemfile Gemfile.lock ./

RUN bundle install

CMD ["ruby", "importer.rb"]
