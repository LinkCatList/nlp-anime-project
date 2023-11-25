# nlp-anime-project

Цель - создать нейронную сеть, которая будет по описанию аниме определять его возможный жанр.

<img src="https://github.com/zavman58/nlp-anime-project/blob/main/pic/tyan.png" width="200" height="250"> <img src="https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/tyan2.png" width="200" height="250"> <img src="https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/tyan3.png" width="200" height="250"> 


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
Посчитаем среднюю длину описания, чтобы определить длнну последовательности

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

Для каждого датасета применим pad_sequences, для того чтобы длины тензоров были одиковые (иначе модель будет жаловаться). Дополним их незначащими нулями в конце

```python
MAX_SEQ_LEN = 70

train_data = tf.keras.preprocessing.sequence.pad_sequences(
    train_data,
    value= vocab['<PAD>'],
    padding= 'post',
    maxlen= MAX_SEQ_LEN)

test_data = tf.keras.preprocessing.sequence.pad_sequences(
    test_data,
    value= vocab['<PAD>'],
    padding= 'post',
    maxlen= MAX_SEQ_LEN)

>>>Тренировочные данные:
[   4  678  580    2    2   23    2   40    5   89  296    2   10 1689
  347    2 1562 4625    2    9 2074    5 3090 5733    8 4321    2   82
    3  879    6    3    2   60   56   35    2    3   57  613    3 2342
   51    8 1343    9   17   25   48 3002    2  281 1166  448   17 2854
 4625    7   10  160  365    4    3  281   82    5   44 2534 1724   55]

   Тестовые данные:
[   1  218   21 1110 1026   22    2    3 8804   27    2   30   40    8
    5  942    6 2891    2    8 1273 1301   14    2    9  870 6425    8
    3 1110   72 2503    7   40   54   19    3  238   11   76    0    0
    0    0    0    0    0    0    0    0    0    0    0    0    0    0
    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
```

## Обучаем модель

Разбиваем на тренировочную и тестовую выборки 

```python
partial_x_train, x_val, partial_y_train, y_val = train_test_split(train_data, train_label,
                                                                  test_size = 0.10, random_state = 42)
```

Создадим объект model и пропишем архитектуру

```python
EMB_SIZE = 32

model = tf.keras.Sequential([
    tf.keras.layers.Embedding(VOCAB_SIZE, EMB_SIZE),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.1, recurrent_dropout=0.1)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.2, recurrent_dropout=0.1)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=False, dropout=0.2, recurrent_dropout=0.1)),
    tf.keras.layers.Dense(CLASS_NUM, activation= 'softmax'),
])



>>>
_________________________________________________________________
 Layer (type)                Output Shape              Param #   
=================================================================
 embedding (Embedding)       (None, None, 32)          320000    
                                                                 
 bidirectional (Bidirection  (None, None, 64)          16640     
 al)                                                             
                                                                 
 bidirectional_1 (Bidirecti  (None, None, 64)          24832     
 onal)                                                           
                                                                 
 bidirectional_2 (Bidirecti  (None, 64)                24832     
 onal)                                                           
                                                                 
 dense (Dense)               (None, 10)                650       
                                                                 
=================================================================
Total params: 386954 (1.48 MB)
Trainable params: 386954 (1.48 MB)
Non-trainable params: 0 (0.00 Byte)
_________________________________________________________________
```


Приступаем к обучению сетки

```python
BATCH_SIZE = 64
NUM_EPOCHS = 30

cpt_path = 'data/14_text_classifier.hdf5'
checkpoint = tf.keras.callbacks.ModelCheckpoint(cpt_path, monitor='acc', verbose=1, save_best_only= True, mode='max')

model.compile(loss='categorical_crossentropy', optimizer='adam', metrics=['acc'])

history = model.fit(partial_x_train, partial_y_train, validation_data= (x_val, y_val),
                   epochs= NUM_EPOCHS, batch_size= BATCH_SIZE, verbose= 1,
                   callbacks=[checkpoint])
```

Ждемс пару часиков и смотрим на result:

```python
results = model.evaluate(test_data, test_label)
```


![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/graf6.png)

![alt text](https://github.com/LinkCatList/nlp-anime-project/blob/main/pic/graph7.png)

78% точности


Создадим функцию, которая вовращает предсказание по промту

```python
def get_result_model(prompt):
    prompt_tokenize = tokenize_text(prompt)

    prompt_extra = []
    for i in prompt_tokenize:
      if i in mini_vocab:
        prompt_extra.append(i)

    prompt_encode = EncodeTokens(prompt_extra, mini_vocab)
    prediction = my_model.predict(tf.constant([prompt_encode]))

    return labels[np.argmax(prediction)]
```

Для примера работы модели возмем аниме "Святой воин Цербер". На jut.su он определяется жанром "приключения" и "фэнтази".
С того же сайта возьмем описание и дадим в качестве запроса модели (не забываем переводить на английский). 

Смотрим на результат:
```python
text = 'Dangerous adventures, massive battles, as well as the use of ancient magic can be seen in the anime "Seisen Cerberus: Ryuukoku no Fatalites". Ahead of Hiro are waiting for difficult training, as well as the search for answers to questions that are related to the sworn enemy. A skilled swordsman together with his companions will meet many enemies and friends who will influence their destinies.'

get_result_model(text)

>>> 'Adventure'

```
Наша модель молодец, она верно определила жанр аниме

Потыкать нейроночку здесь:
https://clck.ru/36mVwp

Всем さようなら!




