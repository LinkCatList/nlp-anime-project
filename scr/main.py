"""
Установим библиотеки
"""

import pandas as pd
import numpy as np
import nltk
import string
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
import tensorflow as tf
nltk.download('punkt')

from keras.models import load_model
from googletrans import Translator

"""#  Предобработка данных

В этом разделе мы подготовим наши данные к токенизации и прочим штукам

Train dataset:
"""

# загрузка
train_df = pd.read_csv('/content/train_new.csv')
train_df = train_df.rename(columns={'Unnamed: 0': 'Id'})

# создаем датафрейм с кол-во жанров в трейне
genre_count_df = train_df[['Id','genres']].groupby('genres').agg('count').sort_values('Id', ascending = False)

genre_count_df.reset_index(inplace = True)
genre_count_df.rename(columns ={'Id':'genre_cnt'}, inplace = True)

# рисуем график
genre_count_df.plot.bar(y = 'genre_cnt', x = 'genres')

# объединяем
train_df = train_df.merge(genre_count_df, how = 'left', left_on='genres', right_on='genres')

# убираем жанры, кол-во которых меньше min_cnt
min_cnt = 700
train_df = train_df[train_df['genre_cnt'] > min_cnt]
train_df['genre_cnt'] = train_df['genre_cnt'].replace('nan', 0)

train_df.genres.value_counts() # получаем 10 жанров, которые позже будем предсказывать

"""Test dataset:"""

# загрузка
test_df =  pd.read_csv('/content/test_new.csv')
test_df = test_df.rename(columns={'Unnamed: 0': 'Id'})

# создаем датафрейм с кол-во жанров в тесте
genre_count_df_test = test_df[['Id','genres']].groupby('genres').agg('count').sort_values('Id', ascending = False)

genre_count_df_test.reset_index(inplace = True)
genre_count_df_test.rename(columns ={'Id':'genre_cnt'}, inplace = True)

# рисуем график
genre_count_df_test.plot.bar(y = 'genre_cnt', x = 'genres')

# объединяем
test_df = test_df.merge(genre_count_df_test, how = 'left', left_on='genres', right_on='genres')

# убираем жанры, кол-во которых меньше min_cnt
min_cnt = 700
test_df = test_df[test_df['genre_cnt'] > min_cnt]
test_df['genre_cnt'] = test_df['genre_cnt'].replace('nan', 0)

# уберем пропуски в столбцах description и genres

train_df = train_df.dropna(subset=['description', 'genres'])
test_df = test_df.dropna(subset=['description', 'genres'])

"""# Токенизация"""

# создадим функцию, которая будет токенизировать текст
def tokenize_text(text):
  tokenized_text = nltk.word_tokenize(text)
  tokens = [i.lower() for i in tokenized_text if ( i not in string.punctuation )]
  return tokens

# применяем к нашим данным
train_df['description_tokenized'] = train_df['description'].apply(lambda x:tokenize_text(x))
test_df['description_tokenized'] = test_df['description'].apply(lambda x:tokenize_text(x))

"""# Создание словаря"""

vocab = {}
mx = 1000000

vocab["<PAD>"] = mx + 2;
vocab["<START>"] = mx + 1;
vocab["<UNKW>"] = mx;

for tokens in train_df['description_tokenized']:
    for word in tokens:
        if word not in vocab.keys():
            vocab[word] = 1
        else:
            vocab[word] += 1

vocab = {k: v for k, v in sorted(vocab.items(), key = lambda item: item[1], reverse = True)}

cnt = 0
for k in vocab.keys():
  vocab[k] = cnt
  cnt += 1

print(len(vocab))
# print(vocab)

"""Удаляем слова, которые слишком редко встречались, это будет основной словарь для обучения"""

mini_vocab = {}

for k, v in vocab.items():
  if v < 10000:
    mini_vocab[k] = v;

print(len(mini_vocab))
# print(mini_vocab)

"""# Энкодер и Декодер"""

def EncodeTokens(ListTokens, VocabTokens): # токены -> числа
  res = []
  res = [VocabTokens.get(word, VocabTokens['<UNKW>']) for word in ListTokens]
  return [VocabTokens['<START>']] + res

def DecodeTokens(EncodedTokens, VocabTokens): # числа -> токены
  res = []
  for i in EncodedTokens:
    for word, ind in VocabTokens.items():
      if i == ind:
        res.append(word)
        break
  return res

# применяем к данным

train_df['description_encoded'] = train_df['description_tokenized'].apply(lambda x: EncodeTokens(x, mini_vocab))
test_df['description_encoded'] = test_df['description_tokenized'].apply(lambda x: EncodeTokens(x, mini_vocab))

"""Переведем в нужные нам типы данных, посчитаем описательные статистики и выведем график"""

train_data = train_df.description_encoded.to_numpy()
train_label = pd.get_dummies(train_df['genres']).values

test_data = test_df.description_encoded.to_numpy()
test_label = pd.get_dummies(test_df['genres']).values

train_df['description_len'] = train_df['description_encoded'].apply (len)

print ('минимальная длина описания:', train_df.description_len.min())
print ('средняя длина описания:', round(train_df.description_len.mean()))
print ('максимальная длина описания:', train_df.description_len.max())

plt.hist(train_df.description_len, density = True)

"""# Обучение

Проведем preprocessing данных. Мы добиваем неважными нулями тензоры, потому что иначе модель это не скушает
"""

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

#print('Тернировочные данные:')
#print(train_data.shape)
#print(train_data[0])
#print()
#print('Тестовые данные:')
#print(test_data.shape)
#print(test_data[0])

"""Разбиваем на тренировочную и тестовые выборки"""

partial_x_train, x_val, partial_y_train, y_val = train_test_split(train_data, train_label,
                                                                  test_size = 0.10, random_state = 42)

"""Создаем объект класса и прописываем архитектуру"""

VOCAB_SIZE = len(mini_vocab)
EMB_SIZE = 32
CLASS_NUM = y_val.shape[1]

model = tf.keras.Sequential([
    tf.keras.layers.Embedding(VOCAB_SIZE, EMB_SIZE),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.1, recurrent_dropout=0.1)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.2, recurrent_dropout=0.1)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=True, dropout=0.2, recurrent_dropout=0.1)),
    tf.keras.layers.Bidirectional(
        tf.keras.layers.LSTM(EMB_SIZE, return_sequences=False, dropout=0.2, recurrent_dropout=0.1)),
    tf.keras.layers.Dense(CLASS_NUM, activation= 'softmax'),
])

model.summary()

"""Начинаем самое обучение модели"""

BATCH_SIZE = 64
NUM_EPOCHS = 30

# сохраняем файл с весами и параметрами модели
cpt_path = '/content/model_history.hdf5'
checkpoint = tf.keras.callbacks.ModelCheckpoint(cpt_path, monitor='acc', verbose=1, save_best_only= True, mode='max')

model.compile(loss='categorical_crossentropy', optimizer='adam', metrics=['acc'])

# history = model.fit(partial_x_train, partial_y_train, validation_data= (x_val, y_val),
                   #epochs= NUM_EPOCHS, batch_size= BATCH_SIZE, verbose= 1,
                   #callbacks=[checkpoint])

"""Two hour later...

Выводим хорошие (нет) графики loss и accuracy
"""

epochs = range(1, len(history.history['acc']) + 1)

plt.figure()
plt.plot(epochs, history.history['loss'], 'bo', label='Training loss')
plt.plot(epochs, history.history['val_loss'], 'b', label='Validation loss')
plt.title('Training and validation loss')
plt.xlabel('Epochs')
plt.ylabel('Loss')
plt.legend()
plt.grid()

plt.figure()
plt.plot(epochs, history.history['acc'], 'bo', label='Training mse')
plt.plot(epochs, history.history['val_acc'], 'b', label='Validation mse')
plt.title('Training and validation acc')
plt.xlabel('Epochs')
plt.ylabel('acc')
plt.legend()
plt.grid()

"""# Тестирование"""

# наши 10 классов
labels = train_df.genres.unique()

# модель, что получилась после обучения
my_model = load_model('/content/14_text_classifier.hdf5')

# вводим наш запрос - описание аниме, у которого хотим определить жанр
prompt = '12345678'

# токенизируем
prompt_tokenize = tokenize_text(prompt)

# удаляем лишние слова, которые не содержаться в словаре mini_vocab
prompt_extra = []
for i in prompt_tokenize:
  if i in mini_vocab:
    prompt_extra.append(i)

# энкодим (переводим в числа)
prompt_encode = EncodeTokens(prompt_extra, mini_vocab)

# получаем предсказание в числах с помощью модели
prediction = my_model.predict(tf.constant([prompt_encode]))

# переводим на человеческий язык
ans = labels[np.argmax(prediction)]
ans
