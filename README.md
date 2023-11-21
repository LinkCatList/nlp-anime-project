# nlp-anime-project

Цель - создать нейронную сеть, которая будет по описанию аниме определять его возможный жанр.

<img src="https://github.com/zavman58/nlp-anime-project/blob/main/pic/tyan.png" width="200" height="250"> <img src="https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/tyan2.png" width="200" height="250"> <img src="https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/tyan3.png" width="200" height="250"> <img src="https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/tyan4.jpg" width="250" height="250">


## Загрузка и первоначальная обработка данных
Загружаем данные для обучения по [ссылке](https://github.com/LinkCatList/nlp-anime-project/blob/main/datasets/anime-description.csv)  
   
Для каждого аниме сохраним два основых жанра
```python
train_df['genres'] = train_df['genres'].str.split(', ')
train_df['first_genre'] = train_df['genres'].apply(lambda x: x[0])
train_df['second_genre'] = train_df['genres'].apply(lambda x: x[1] if len(x) > 1 else x[0])
```
Посмотрим на графики распределения первого и второго жанров:

![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/graph.png)

![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/graph2.png)

Для каждого аниме посчитаем сколько раз всего встречается первый жанр + сколько всего раз встречается второй жанр и удалим те, жанры которых встречаются редко

```python
min_cnt = 1000
train_df = train_df[train_df['genre_cnt'] > min_cnt]
```

## Подготовка данных для модели
Напишем функцию токенизации текста и разобьем каждое описание на токены.
```python
def tokenize_text(text):
  tokenized_text = nltk.word_tokenize(text)
  tokens = [i.lower() for i in tokenized_text if ( i not in string.punctuation )]
  return tokens
```

```python
train_df['description_tokenized'] = train_df['description'].apply(lambda x:tokenize_text(x))
test_df['description_tokenized'] = test_df['description'].apply(lambda x:tokenize_text(x))
```

Теперь для каждого слова посчитаем количество вхождений во все описания:
```python
for tokens in train_df['description_tokenized']:
  for word in tokens:
    if word not in vocab.keys():
      vocab[word] = 1
    else:
      vocab[word]+=1

```
Получится большой-большой словарь, но для обучения можно использовать облегченный словарь, возьмем слова которые встречаются не так часто:
```python
mini_vocab = {}
for k, v in vocab.items():
  if v<10000:
    mini_vocab[k] = v

```

Напишем функцию кодировки и декодировки описаний.      
Закодируем все описания и добавим их в датафрейм:

![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/table1.png)

## Готовим данные к обучению
Посчитаем среднюю длинну описания, чтобы определить длинну последовательности

```python
train_df['description_len'] = train_df['description_encoded'].apply (len)

print ('минимальная длина описания:', train_df.description_len.min())
print ('средняя длина описания:', round(train_df.description_len.mean()))
print ('максимальная длина описания:', train_df.description_len.max())

plt.hist(train_df.description_len, density = True)

>>> минимальная длина описания: 3
    средняя длина описания: 92
    максимальная длина описания: 487
```

![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/graph3.png)
